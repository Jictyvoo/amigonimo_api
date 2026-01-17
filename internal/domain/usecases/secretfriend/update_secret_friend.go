package secretfriend

import (
	"fmt"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type UpdateInput struct {
	ID       entities.HexID
	Name     *string
	Datetime *time.Time
	Location *string
}

func (uc *UseCase) Update(input UpdateInput) error {
	sf, err := uc.repo.GetSecretFriendByID(input.ID)
	if err != nil {
		return fmt.Errorf("get for update: %w", err)
	}

	if input.Name != nil {
		sf.Name = *input.Name
	}
	if input.Datetime != nil {
		sf.Datetime = input.Datetime
	}
	if input.Location != nil {
		sf.Location = *input.Location
	}
	sf.UpdatedAt = time.Now()

	if err = uc.repo.UpdateSecretFriend(&sf); err != nil {
		return fmt.Errorf("update secret friend: %w", err)
	}

	return nil
}
