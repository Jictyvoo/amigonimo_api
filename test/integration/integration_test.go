package integration

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy"
	"github.com/wrapped-owls/testereiro/puppetest"

	"github.com/jictyvoo/amigonimo_api/build/migrations"
	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
)

var NewEngine func(t testing.TB) *puppetest.Engine

func mysqlPerformer(ctx context.Context, conf puppetest.DBConnectionConfig) (*sql.DB, error) {
	appConf := bootstrap.Config()
	dbConf := appConf.Database
	dbConf.Database = ""
	if conf.DBName != "" {
		dbConf.Database = conf.DBName
	}

	mysqlConf := bootstrap.MySQLConfig(dbConf)
	mysqlConf.MultiStatements = conf.AllowMultiStatements

	return bootstrap.PerformDatabaseConnection(ctx, mysqlConf, dbConf.Timeout)
}

func withTestServer(e *puppetest.Engine) (http.Handler, error) {
	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	remy.RegisterInstance(inj, e.DB())

	conf := bootstrap.Config()
	bootstrap.DoInjections(inj, conf)

	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	remy.RegisterInstance(inj, privateKey)
	server, servErr := bootstrap.NewWebServer(conf, &privateKey.PublicKey, inj)
	if servErr != nil {
		return nil, fmt.Errorf("failed to create web server: %w", servErr)
	}

	return server.Mux, nil
}

func TestMain(m *testing.M) {
	engineFactory, err := puppetest.NewEngineFactory(
		puppetest.WithConnectionFactory(mysqlPerformer, true),
		puppetest.WithExtensions(
			puppetest.WithMigrationRunner(migrations.VersionedMigrationsFS()),
			puppetest.WithTestServerFromEngine(withTestServer),
		),
	)
	if err != nil {
		slog.Error("failed to setup engine factory", slog.String("error", err.Error()))
		os.Exit(1)
	}

	NewEngine = engineFactory.NewEngine
	code := m.Run() // Run all tests in the package

	// Teardown/cleanup
	if err = engineFactory.Close(); err != nil {
		slog.Error("failed to close engine factory", slog.String("error", err.Error()))
		os.Exit(1)
	}

	os.Exit(code)
}
