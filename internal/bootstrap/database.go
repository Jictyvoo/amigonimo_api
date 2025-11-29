package bootstrap

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func OpenDatabase(dbConfig config.Database) *sql.DB {
	// Build MySQL DSN: user:password@tcp(host:port)/database?parseTime=true&loc=UTC
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC",
		dbConfig.User, dbConfig.Password,
		dbConfig.Host, dbConfig.Port, dbConfig.Database,
	)

	db, dbErr := sql.Open("mysql", dsn)
	if dbErr != nil {
		slog.Error("failed to open database", slog.String("error", dbErr.Error()))
		os.Exit(1)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return db
}
