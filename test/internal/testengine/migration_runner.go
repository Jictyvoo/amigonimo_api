package testengine

import (
	"database/sql"
	"io/fs"
	"strings"
	"testing"

	"github.com/jictyvoo/amigonimo_api/build/migrations"
)

func WithMigrationRunner(t testing.TB, engine *Engine) {
	runMigrations(t, engine.db)
}

func runMigrations(t testing.TB, db *sql.DB) {
	migrationFS := migrations.VersionedMigrationsFS()
	entries, err := fs.ReadDir(migrationFS, ".")
	if err != nil {
		t.Fatalf("error reading migrations: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		var content []byte
		content, err = fs.ReadFile(migrationFS, entry.Name())
		if err != nil {
			t.Fatalf("failed to read migration file `%s`: %v", entry.Name(), err.Error())
		}

		_, err = db.ExecContext(t.Context(), string(content))
		if err != nil {
			t.Fatalf("failed to execute migration file `%s`: %v", entry.Name(), err.Error())
		}
	}
}
