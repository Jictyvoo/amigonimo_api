package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// ParticipantFacade defines what denylist needs from participant.
type ParticipantFacade interface {
	ports.Facade
	CheckParticipantInSecretFriend(sfID, userID entities.HexID) (entities.Participant, error)
}

// SecretFriendFacade defines what denylist needs from secretfriend.
type SecretFriendFacade interface {
	ports.Facade
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
}

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
