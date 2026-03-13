package secretfriend

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=secretfriend github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend Repository

type Repository interface {
	CreateSecretFriend(sf *entities.SecretFriend) error
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateSecretFriend(sf *entities.SecretFriend) error
	ListSecretFriends(userID entities.HexID) ([]entities.SecretFriend, error)
	GetSecretFriendByInviteCode(code string) (entities.SecretFriend, error)
}

type UseCase struct {
	repo           Repository
	associatedUser entities.User
}

func New(associatedUser entities.User, repo Repository) UseCase {
	return UseCase{associatedUser: associatedUser, repo: repo}
}

func (uc *UseCase) Get(id entities.HexID) (entities.SecretFriend, error) {
	sf, err := uc.repo.GetSecretFriendByID(id)
	if err != nil {
		return entities.SecretFriend{}, fmt.Errorf("get secret friend: %w", err)
	}
	return sf, nil
}

func (uc *UseCase) CheckUserIsOwner(sfID entities.HexID) (bool, error) {
	sf, err := uc.repo.GetSecretFriendByID(sfID)
	if err != nil {
		return false, fmt.Errorf("get secret friend: %w", err)
	}
	return sf.OwnerID == uc.associatedUser.ID, nil
}
