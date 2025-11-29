package domain

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services"
)

func RegisterServices(inj remy.Injector) {
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[*services.AuthService], services.NewAuthService,
	)
}
