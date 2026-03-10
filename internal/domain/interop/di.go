package interop

import (
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/facades"
)

func RegisterFacades(inj remy.Injector) {
	remy.RegisterConstructorArgs1(
		inj, remy.Factory[*facades.SecretFriendFacade], facades.NewSecretFriendFacade,
	)

	remy.RegisterConstructorArgs1(
		inj, remy.Factory[*facades.ParticipantFacade], facades.NewParticipantFacade,
	)
}
