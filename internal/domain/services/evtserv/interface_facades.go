package evtserv

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate mockgen -destination=../mocks/secretfriend_repo_mock.go -package=mocks github.com/jictyvoo/amigonimo_api/internal/domain/secretfriendserv SecretFriendRepository

type SecretFriendRepository interface {
	CreateSecretFriend(secretFriend *entities.SecretFriend) error
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateSecretFriend(secretFriend *entities.SecretFriend) error
	GetParticipantsCount(secretFriendID entities.HexID) (int, error)
	GetDrawResultForUser(secretFriendID, userID entities.HexID) (entities.DrawResultItem, error)
}
