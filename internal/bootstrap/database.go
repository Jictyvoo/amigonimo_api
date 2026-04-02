package bootstrap

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func OpenDatabase(dbConfig config.Database) *sql.DB {
	mysqlConf := mysqlrepo.MySQLConfig(dbConfig)
	db, err := mysqlrepo.PerformDatabaseConnection(
		context.Background(), mysqlConf, mysqlConf.Timeout,
	)
	if err != nil {
		slog.Error("failed to start *sql.DB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return db
}
