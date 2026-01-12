package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("failed to ping database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return db
}
