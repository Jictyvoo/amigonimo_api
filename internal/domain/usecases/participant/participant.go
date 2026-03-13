package participant

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=secret_friend_facade_mock_test.go -package=participant github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant SecretFriendFacade

// SecretFriendFacade defines what participant module needs from secretfriend.
type SecretFriendFacade interface {
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	CheckUserIsOwner(sfID entities.HexID) (bool, error)
}

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=participant github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant Repository

type Repository interface {
	AddParticipant(sfID, userID entities.HexID) (entities.Participant, error)
	ListParticipants(sfID entities.HexID) ([]entities.Participant, error)
	SetParticipantReady(sfID, userID entities.HexID, isReady bool) error
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
