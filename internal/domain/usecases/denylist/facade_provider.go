package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=facade_provider_mock_test.go -package=denylist github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist participantFacadePort,secretFriendFacadePort

// ParticipantFacade defines what denylist needs from participant.
type (
	participantFacadePort interface {
		CheckParticipantInSecretFriend(sfID, userID entities.HexID) (entities.Participant, error)
	}
	ParticipantFacade interface {
		ports.Facade
		participantFacadePort
	}
)

type (
	secretFriendFacadePort interface {
		GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	}
	// SecretFriendFacade defines what denylist needs from secretfriend.
	SecretFriendFacade interface {
		ports.Facade
		secretFriendFacadePort
	}
)

type FacadeProvider struct {
	participant  ParticipantFacade
	secretFriend SecretFriendFacade
}

func NewFacadeProvider(
	validator ParticipantFacade,
	sfProvider SecretFriendFacade,
) FacadeProvider {
	return FacadeProvider{
		participant:  validator,
		secretFriend: sfProvider,
	}
}
