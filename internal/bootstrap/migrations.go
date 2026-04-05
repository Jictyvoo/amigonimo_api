package bootstrap

import (
	"context"
	"io/fs"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/migratorrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

// RunMigrations opens a dedicated MultiStatements=true connection, applies
// all pending SQL migrations from migrationsFS, then closes the connection.
func RunMigrations(dbConf config.Database, migrationsFS fs.ReadDirFS) error {
	mysqlConf := mysqlrepo.MySQLConfig(dbConf)
	mysqlConf.MultiStatements = true

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := mysqlrepo.PerformDatabaseConnection(ctx, mysqlConf, dbConf.Timeout)
	if err != nil {
		return err
	}
	defer db.Close()

	return migratorrepo.New(db).Run(migrationsFS)
}
