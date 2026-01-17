package secretfriend

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Repository interface {
	CreateSecretFriend(sf *entities.SecretFriend) error
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateSecretFriend(sf *entities.SecretFriend) error
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Get(id entities.HexID) (entities.SecretFriend, error) {
	sf, err := uc.repo.GetSecretFriendByID(id)
	if err != nil {
		return entities.SecretFriend{}, fmt.Errorf("get secret friend: %w", err)
	}
	return sf, nil
}
