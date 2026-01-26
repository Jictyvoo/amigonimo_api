package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func PerformDatabaseConnection(
	ctx context.Context, conf mysql.Config, pingTimeout time.Duration,
) (*sql.DB, error) {
	dsn := conf.FormatDSN()
	db, dbErr := sql.Open("mysql", dsn)
	if dbErr != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", dbErr)
	}

	if pingTimeout <= 0 {
		pingTimeout = time.Second
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func OpenDatabase(dbConfig config.Database) *sql.DB {
	mysqlConf := MySQLConfig(dbConfig)
	db, err := PerformDatabaseConnection(context.Background(), mysqlConf, mysqlConf.Timeout)
	if err != nil {
		slog.Error("failed to start *sql.DB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return db
}

func MySQLConfig(dbConf config.Database) mysql.Config {
	return mysql.Config{
		User:                 dbConf.User,
		Passwd:               dbConf.Password,
		Net:                  "tcp",
		Addr:                 dbConf.Host + ":" + strconv.FormatUint(uint64(dbConf.Port), 10),
		DBName:               dbConf.Database,
		MultiStatements:      false,
		ParseTime:            true,
		AllowNativePasswords: true,
		Loc:                  time.UTC,
	}
}
