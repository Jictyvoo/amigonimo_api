package sqltest

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jictyvoo/amigonimo_api/build/migrations"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

// Factory manages a root DB connection used to create an isolated database for
// each test. Each NewDB call creates a fresh database, applies migrations, and
// registers a cleanup that drops the database when the test finishes.
type Factory struct {
	rootDB *sql.DB
	dbConf config.Database
}

// NewFactory loads configuration, then opens a root connection (no database
// selected) that has enough privileges to CREATE/DROP databases.
func NewFactory() (*Factory, error) {
	conf, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("sqltest: load config: %w", err)
	}

	rootConf := conf.Database
	rootConf.Database = "" // connect without selecting a database

	mysqlConf := mysqlrepo.MySQLConfig(rootConf)
	mysqlConf.MultiStatements = true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rootDB, err := mysqlrepo.PerformDatabaseConnection(ctx, mysqlConf, conf.Database.Timeout)
	if err != nil {
		return nil, fmt.Errorf("sqltest: open root connection: %w", err)
	}

	return &Factory{rootDB: rootDB, dbConf: conf.Database}, nil
}

// NewDB creates a fresh database named after t.Name(), applies all migrations,
// and returns a connection to it. The database is dropped automatically when
// the test finishes via t.Cleanup.
func (f *Factory) NewDB(t testing.TB) *sql.DB {
	t.Helper()

	dbName := dbNameNormalizer(t.Name())
	ctx := t.Context()

	if _, err := f.rootDB.ExecContext(
		ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName),
	); err != nil {
		t.Fatalf("sqltest: create database %q: %v", dbName, err)
	}

	dbConf := f.dbConf
	dbConf.Database = dbName

	mysqlConf := mysqlrepo.MySQLConfig(dbConf)
	mysqlConf.MultiStatements = true

	db, err := mysqlrepo.PerformDatabaseConnection(ctx, mysqlConf, f.dbConf.Timeout)
	if err != nil {
		t.Fatalf("sqltest: connect to %q: %v", dbName, err)
	}

	if err = applyMigrations(db, migrations.VersionedMigrationsFS()); err != nil {
		_ = db.Close()
		t.Fatalf("sqltest: apply migrations to %q: %v", dbName, err)
	}

	t.Cleanup(
		func() {
			if _, cleanErr := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)); cleanErr != nil {
				t.Errorf("failed to exec database drop: %s", cleanErr.Error())
			}
			if cleanErr := db.Close(); cleanErr != nil {
				t.Errorf("failed to close connection to test database: %s", cleanErr.Error())
			}
		},
	)

	return db
}

// Close shuts down the root connection. Call in TestMain after m.Run().
func (f *Factory) Close() error {
	if f.rootDB != nil {
		err := f.rootDB.Close()
		f.rootDB = nil
		return err
	}
	return nil
}
