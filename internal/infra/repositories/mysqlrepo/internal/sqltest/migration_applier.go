package sqltest

import (
	"database/sql"
	"fmt"
	"io/fs"
	"strings"
)

func applyMigrations(db *sql.DB, migrationsFS fs.ReadDirFS) error {
	entries, err := migrationsFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, fileEntry := range entries {
		if fileEntry.IsDir() || !strings.HasSuffix(fileEntry.Name(), ".sql") {
			continue
		}

		fileName := fileEntry.Name()
		data, readErr := fs.ReadFile(migrationsFS, fileName)
		if readErr != nil {
			return fmt.Errorf("read migration %q: %w", fileName, readErr)
		}
		if _, execErr := db.Exec(string(data)); execErr != nil {
			return fmt.Errorf("exec migration %q: %w", fileName, execErr)
		}
	}
	return nil
}
