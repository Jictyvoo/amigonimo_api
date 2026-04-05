package migratorrepo

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"
	"strings"
	"time"
)

// SchemaRevision mirrors the schema_revisions table row.
type SchemaRevision struct {
	Version       string
	Description   string
	Applied       int64
	Total         int64
	ExecutedAt    *time.Time
	ExecutionTime int64 // milliseconds
	Error         *string
	Hash          string
}

// Migrator applies SQL migrations from an embedded FS and tracks them in
// the schema_revisions table (structurally similar to atlas_schema_revisions).
type Migrator struct {
	db *sql.DB
}

func New(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Run ensures the tracking table exists, seeds from atlas_schema_revisions
// if present, then applies every pending migration file in order.
// The caller is responsible for providing a db opened with MultiStatements=true.
func (m *Migrator) Run(migrationsFS fs.ReadDirFS) error {
	if _, err := m.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("create schema_revisions: %w", err)
	}

	if err := seedFromAtlas(m.db); err != nil {
		return fmt.Errorf("seed from atlas: %w", err)
	}

	applied, err := loadApplied(m.db)
	if err != nil {
		return fmt.Errorf("load applied revisions: %w", err)
	}

	entries, err := migrationsFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".sql") || applied[name] {
			continue
		}
		if err = apply(m.db, migrationsFS, name); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
		slog.Info("migration applied", slog.String("version", name))
	}
	return nil
}

// seedFromAtlas copies rows from atlas_schema_revisions into schema_revisions
// when the atlas table exists and our table is still empty (first run after switch).
func seedFromAtlas(db *sql.DB) error {
	var atlasExists int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM information_schema.TABLES
 WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'atlas_schema_revisions'`,
	).Scan(&atlasExists)
	if err != nil || atlasExists == 0 {
		return nil
	}

	var count int
	if err = db.QueryRow(`SELECT COUNT(*) FROM schema_revisions`).Scan(&count); err != nil ||
		count > 0 {
		return nil
	}

	_, seedErr := db.Exec(`
INSERT IGNORE INTO schema_revisions (version, description, applied, total, executed_at, execution_time, hash)
SELECT version, description, applied, total, executed_at, execution_time, hash
FROM   atlas_schema_revisions`)
	if seedErr != nil {
		slog.Warn(
			"could not seed from atlas_schema_revisions",
			slog.String("error", seedErr.Error()),
		)
	}
	return nil
}

// loadApplied returns the set of versions already recorded in schema_revisions.
func loadApplied(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT version FROM schema_revisions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err = rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, rows.Err()
}

// apply executes a single SQL migration file and records the revision.
func apply(db *sql.DB, migrationsFS fs.ReadDirFS, name string) error {
	data, err := fs.ReadFile(migrationsFS, name)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	start := time.Now()

	_, execErr := db.Exec(string(data))

	execMs := time.Since(start).Milliseconds()

	if execErr != nil {
		errMsg := execErr.Error()
		_ = recordRevision(db, SchemaRevision{
			Version:       name,
			Description:   name,
			Applied:       0,
			Total:         1,
			ExecutionTime: execMs,
			Error:         &errMsg,
			Hash:          hash,
		})
		return execErr
	}

	now := time.Now()
	return recordRevision(db, SchemaRevision{
		Version:       name,
		Description:   name,
		Applied:       1,
		Total:         1,
		ExecutedAt:    &now,
		ExecutionTime: execMs,
		Hash:          hash,
	})
}

func recordRevision(db *sql.DB, rev SchemaRevision) error {
	_, err := db.Exec(
		`INSERT INTO schema_revisions
 (version, description, applied, total, executed_at, execution_time, error, hash)
 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rev.Version, rev.Description, rev.Applied, rev.Total,
		rev.ExecutedAt, rev.ExecutionTime, rev.Error, rev.Hash,
	)
	return err
}
