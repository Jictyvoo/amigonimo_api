package migratorrepo

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"slices"
	"strings"
	"time"
)

// Migrator applies SQL migrations from an embedded FS and tracks them in the
// schema_revisions table. The provided db must be opened with MultiStatements=true.
type Migrator struct {
	db *sql.DB
}

func New(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Run creates the tracking table if needed, then applies every pending
// migration file from migrationsFS in lexicographic (chronological) order.
func (m *Migrator) Run(migrationsFS fs.ReadDirFS) error {
	if _, err := m.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("create schema_revisions: %w", err)
	}

	applied, err := loadApplied(m.db)
	if err != nil {
		return fmt.Errorf("load applied revisions: %w", err)
	}

	entries, err := migrationsFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	slices.SortFunc(
		entries, func(a, b fs.DirEntry) int {
			return strings.Compare(a.Name(), b.Name())
		},
	)

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".sql") || applied[name] {
			continue
		}
		if err = m.apply(migrationsFS, name); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
		slog.Info("migration applied", slog.String("version", name))
	}
	return nil
}

// apply executes a single migration file as one multi-statement query and
// records the revision. Statement count is derived from the file content for
// accurate applied/total tracking.
func (m *Migrator) apply(migrationsFS fs.ReadDirFS, name string) error {
	data, err := fs.ReadFile(migrationsFS, name)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	version, description := parseFilename(name)
	stmts := splitStatements(string(data))
	total := int64(len(stmts))
	start := time.Now()

	_, execErr := m.db.Exec(string(data))

	newRevision := SchemaRevision{
		Version:       version,
		Description:   description,
		Applied:       total,
		Total:         total,
		ExecutionTime: time.Since(start).Milliseconds(),
		Hash:          fmt.Sprintf("%x", sha256.Sum256(data)),
	}
	if execErr != nil {
		newRevision.Error = new(execErr.Error())
		newRevision.Applied = 0
	}

	return recordRevision(m.db, newRevision)
}
