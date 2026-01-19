package bootstrap

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func DoInjections(inj remy.Injector, conf config.Config) {
	remy.RegisterInstance(inj, conf)

	infra.RegisterInfraServices(inj)
	domain.RegisterServices(inj)

	// Helper to retrieve user from context, after token authentication
	remy.RegisterConstructorArgs1Err(inj, remy.Factory[entities.User], NewUserFromContext)
}
