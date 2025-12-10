package infra

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/infra/facades/mailfacade"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories"
)

func RegisterInfraServices(inj remy.Injector) {
	repositories.RegisterRepositories(inj)
	remy.RegisterConstructor(inj, remy.Singleton[*mailfacade.MailerImpl], mailfacade.NewMailerImpl)
}
