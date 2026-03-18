package secretfriend

import (
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type CreateInput struct {
	Name            string
	Datetime        time.Time
	Location        string
	MaxDenyListSize uint8
}

func (uc *UseCase) Create(input CreateInput) (entities.SecretFriend, error) {
	inviteCode := uuid.New().String()[:8]

	sf := entities.SecretFriend{
		ID:              entities.HexID(uuid.New()),
		Name:            input.Name,
		Datetime:        input.Datetime,
		Location:        input.Location,
		OwnerID:         uc.associatedUser.ID,
		InviteCode:      inviteCode,
		MaxDenyListSize: input.MaxDenyListSize,
		Status:          entities.StatusDraft,
	}
	sf.CreatedAt = time.Now()
	sf.UpdatedAt = sf.CreatedAt

	if err := uc.repo.CreateSecretFriend(&sf); err != nil {
		return entities.SecretFriend{}, apperr.From(
			"secret_friend_create_failed",
			"failed to create secret friend",
			err,
		)
	}

	return sf, nil
}
