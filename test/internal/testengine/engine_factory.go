package testengine

import (
	"context"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
	"github.com/jictyvoo/amigonimo_api/test/internal/dbsetup"
)

type (
	EngineExtension func(t testing.TB, engine *Engine)
	EngineFactory   struct {
		dbFactory  dbsetup.ConnectionFactory
		extensions []EngineExtension
	}
)

func NewEngineFactory(extensions ...EngineExtension) (*EngineFactory, error) {
	conf := bootstrap.Config()
	dbFactory, err := dbsetup.NewConnectionFactory(context.Background(), conf.Database)
	newFactory := &EngineFactory{
		dbFactory:  dbFactory,
		extensions: extensions,
	}

	return newFactory, err
}

func (fac EngineFactory) NewEngine(t testing.TB) *Engine {
	subDb, err := fac.dbFactory.NewDatabase(t.Context(), t.Name())
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{
		db:     subDb.Connection,
		dbName: subDb.Name,
	}

	for _, extension := range fac.extensions {
		extension(t, engine)
	}

	t.Cleanup(
		func() {
			shutdownErr := engine.Teardown()
			if shutdownErr != nil {
				t.Error(shutdownErr)
			}
		},
	)

	return engine
}

func (fac EngineFactory) Close() error {
	return fac.dbFactory.Close()
}
