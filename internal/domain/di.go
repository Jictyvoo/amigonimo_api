package domain

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/services/evtserv"
)

func RegisterServices(inj remy.Injector) {
	remy.RegisterConstructorArgs3(
		inj, remy.Factory[authserv.AuthService], authserv.NewAuthService,
	)

	// General protected services
	remy.RegisterConstructorArgs2(
		inj, remy.Factory[*evtserv.Service], evtserv.NewService,
	)
}
