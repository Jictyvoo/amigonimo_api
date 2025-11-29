package bootstrap

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain"
	"github.com/jictyvoo/amigonimo_api/internal/infra"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func DoInjections(inj remy.Injector, conf config.Config) {
	remy.RegisterInstance(inj, conf)

	infra.RegisterInfraServices(inj)
	domain.RegisterServices(inj)
}
