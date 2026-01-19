package domain

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend"
)

func RegisterServices(inj remy.Injector) {
	remy.RegisterConstructorArgs3(
		inj, remy.Factory[authserv.AuthService], authserv.NewAuthService,
	)

	// General protected services
	remy.RegisterConstructor(inj, remy.Factory[drawserv.Service], drawserv.New)
	remy.RegisterConstructorArgs2(inj, remy.Factory[secretfriend.UseCase], secretfriend.New)
	remy.RegisterConstructorArgs2(inj, remy.Factory[*drawfriends.UseCase], drawfriends.New)
}
