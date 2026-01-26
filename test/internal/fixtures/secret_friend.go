package fixtures

import (
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
)

type SecretFriendBuilder struct {
	instance *genmodels.SecretFriend
}

func NewSecretFriend() *SecretFriendBuilder {
	uid, err := uuid.NewV7()
	if err != nil {
		log.Panicf("failed to generate uuid: %s", err)
	}
	return &SecretFriendBuilder{
		instance: &genmodels.SecretFriend{
			ID:              uid[:],
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Name:            "Secret Friend Group " + uid.String(),
			MaxDenyListSize: 3,
			InviteCode:      uid.String()[:8],
			Status:          string(entities.StatusOpen),
		},
	}
}

func (b *SecretFriendBuilder) WithOwner(owner *genmodels.User) *SecretFriendBuilder {
	b.instance.OwnerID = owner.ID
	return b
}

func (b *SecretFriendBuilder) WithName(name string) *SecretFriendBuilder {
	b.instance.Name = name
	return b
}

func (b *SecretFriendBuilder) Build() *genmodels.SecretFriend {
	return b.instance
}
