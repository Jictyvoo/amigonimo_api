package integration

import (
	"log/slog"
	"os"
	"testing"

	"github.com/jictyvoo/amigonimo_api/test/internal/testengine"
)

var NewEngine func(t testing.TB) *testengine.Engine

func TestMain(m *testing.M) {
	engineFactory, err := testengine.NewEngineFactory(
		testengine.WithTestServer,
		testengine.WithMigrationRunner,
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
