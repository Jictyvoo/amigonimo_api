package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=participant_facade_mock_test.go -package=denylist github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist ParticipantFacade

// ParticipantFacade defines what denylist needs from participant.
type ParticipantFacade interface {
	ports.Facade
	CheckParticipantInSecretFriend(sfID, userID entities.HexID) (entities.Participant, error)
}

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=secret_friend_facade_mock_test.go -package=denylist github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist SecretFriendFacade

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
