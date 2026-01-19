package secretfriend

import (
	"fmt"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type UpdateInput struct {
	ID       entities.HexID
	Name     string
	Datetime time.Time
	Location string
}

func (uc *UseCase) Update(input UpdateInput) error {
	sf, err := uc.repo.GetSecretFriendByID(input.ID)
	if err != nil {
		return fmt.Errorf("get for update: %w", err)
	}

	if sf.OwnerID != uc.associatedUser.ID { // Must ensure that only the owner can change its info
		return fmt.Errorf("not owned by %s", uc.associatedUser.ID)
	}

	if input.Name != "" {
		sf.Name = input.Name
	}
	if !input.Datetime.IsZero() {
		sf.Datetime = input.Datetime
	}
	if input.Location != "" {
		sf.Location = input.Location
	}
	if err = uc.repo.UpdateSecretFriend(&sf); err != nil {
		return fmt.Errorf("update secret friend: %w", err)
	}

	return nil
}
