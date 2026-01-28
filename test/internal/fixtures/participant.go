package fixtures

import (
	"database/sql"
	"log"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
)

type ParticipantBuilder struct {
	instance *genmodels.Participant
}

func NewParticipant() *ParticipantBuilder {
	uid, err := entities.NewHexID()
	if err != nil {
		log.Panicf("failed to generate HexID: %s", err)
	}

	now := time.Now()
	return &ParticipantBuilder{
		instance: &genmodels.Participant{
			ID:        uid[:],
			CreatedAt: now,
			UpdatedAt: now,
			JoinedAt:  sql.NullTime{Time: now, Valid: true},
		},
	}
}

func (b *ParticipantBuilder) WithUser(user *genmodels.User) *ParticipantBuilder {
	b.instance.UserID = user.ID
	return b
}

func (b *ParticipantBuilder) WithSecretFriend(
	secretFriend *genmodels.SecretFriend,
) *ParticipantBuilder {
	b.instance.SecretFriendID = secretFriend.ID
	return b
}

func (b *ParticipantBuilder) Build() *genmodels.Participant {
	return b.instance
}
