package secretfriend

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
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
		return apperr.From(
			"secret_friend_not_found",
			"secret friend not found",
			err,
		)
	}

	if sf.OwnerID != uc.associatedUser.ID { // Must ensure that only the owner can change its info
		return apperr.Forbidden(
			"secret_friend_update_forbidden",
			"you are not allowed to update this secret friend",
			nil,
		)
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
		return apperr.From(
			"secret_friend_update_failed",
			"failed to update secret friend",
			err,
		)
	}

	return nil
}
