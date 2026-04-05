package migratorrepo

import (
	"database/sql"
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

// loadApplied returns the set of versions already recorded in schema_revisions.
func loadApplied(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT version, description FROM schema_revisions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var ver, desc string
		if err = rows.Scan(&ver, &desc); err != nil {
			return nil, err
		}

		migrationName := ver + "_" + desc + ".sql"
		applied[migrationName] = true
	}
	return applied, rows.Err()
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
