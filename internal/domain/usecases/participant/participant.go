package participant

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// SecretFriendFacade defines what participant module needs from secretfriend.
type SecretFriendFacade interface {
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	CheckUserIsOwner(sfID entities.HexID) (bool, error)
}

type Repository interface {
	AddParticipant(sfID, userID entities.HexID) (entities.Participant, error)
	ListParticipants(sfID entities.HexID) ([]entities.Participant, error)
	GetParticipant(sfID, userID entities.HexID) (entities.Participant, error)
	RemoveParticipant(sfID, userID entities.HexID) error
}

type UseCase struct {
	repo               Repository
	secretFriendFacade SecretFriendFacade
	associatedUser     entities.User
}

func New(
	associatedUser entities.User, repo Repository, sfFacade SecretFriendFacade,
) UseCase {
	return UseCase{
		associatedUser:     associatedUser,
		repo:               repo,
		secretFriendFacade: sfFacade,
	}
}
