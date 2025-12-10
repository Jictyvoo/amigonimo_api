package domain

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
)

func RegisterServices(inj remy.Injector) {
	remy.RegisterConstructorArgs3(
		inj, remy.Factory[authserv.AuthService], authserv.NewAuthService,
	)
}
