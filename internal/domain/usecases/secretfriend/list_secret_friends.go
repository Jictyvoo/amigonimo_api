package secretfriend

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type (
	EventsList struct {
		Created     []entities.SecretFriend
		Participant []entities.SecretFriend
	}
	ActiveInactiveListEvents struct {
		Active, Inactive EventsList
	}
)

func (uc *UseCase) ListUserSecretFriends(userID entities.HexID) (ActiveInactiveListEvents, error) {
	if userID.IsEmpty() {
		userID = uc.associatedUser.ID
	}
	rawList, err := uc.repo.ListSecretFriends(userID)
	if err != nil {
		return ActiveInactiveListEvents{}, apperr.From(
			"secret_friend_list_failed",
			"failed to list secret friends",
			err,
		)
	}

	initialCap := len(rawList) >> 2 //nolint:mnd // divide by 4
	sortedList := ActiveInactiveListEvents{
		Active: EventsList{
			Created:     make([]entities.SecretFriend, 0, initialCap),
			Participant: make([]entities.SecretFriend, 0, initialCap),
		},
		Inactive: EventsList{
			Created:     make([]entities.SecretFriend, 0, initialCap),
			Participant: make([]entities.SecretFriend, 0, initialCap),
		},
	}

	// Iterate over the returned list and sort each status
	for _, sf := range rawList {
		asOwner := !sf.OwnerID.IsEmpty()
		isActive := sf.Status != entities.StatusClosed &&
			// Check for zero datetime, so allow it only if the valid date is after today
			(sf.Datetime.IsZero() || sf.Datetime.After(time.Now()))

		var dest *[]entities.SecretFriend
		switch {
		case asOwner && isActive:
			dest = &sortedList.Active.Created
		case isActive:
			dest = &sortedList.Active.Participant
		case asOwner:
			dest = &sortedList.Inactive.Created
		default:
			dest = &sortedList.Inactive.Participant
		}

		*dest = append(*dest, sf)
	}

	return sortedList, nil
}
