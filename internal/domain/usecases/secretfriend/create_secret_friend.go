package secretfriend

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type CreateInput struct {
	Name            string
	Datetime        *time.Time
	Location        string
	MaxDenyListSize uint8
	OwnerID         entities.HexID
}

func (uc *UseCase) Create(input CreateInput) (entities.SecretFriend, error) {
	inviteCode := uuid.New().String()[:8]

	sf := entities.SecretFriend{
		ID:              entities.HexID(uuid.New()),
		Name:            input.Name,
		Datetime:        input.Datetime,
		Location:        input.Location,
		OwnerID:         input.OwnerID,
		InviteCode:      inviteCode,
		MaxDenyListSize: input.MaxDenyListSize,
		Status:          entities.StatusDraft,
	}
	sf.CreatedAt = time.Now()
	sf.UpdatedAt = sf.CreatedAt

	if err := uc.repo.CreateSecretFriend(&sf); err != nil {
		return entities.SecretFriend{}, fmt.Errorf("create secret friend: %w", err)
	}

	return sf, nil
}
