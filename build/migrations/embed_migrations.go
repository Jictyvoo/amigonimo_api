package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var migrationFiles embed.FS

func VersionedMigrationsFS() fs.ReadDirFS {
	return migrationFiles
}
