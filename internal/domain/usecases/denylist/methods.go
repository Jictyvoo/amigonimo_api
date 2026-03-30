package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc UseCase) GetDenyList(sfID entities.HexID) ([]DeniedEntry, error) {
	deniedUsers, err := uc.repo.GetDenyListByParticipant(
		ParticipantRef{
			UserID:         uc.associatedUser.ID,
			SecretFriendID: sfID,
		},
	)
	if err != nil {
		return nil, apperr.From(
			"denylist_lookup_failed",
			"failed to load denylist",
			err,
		)
	}

	return deniedUsers, nil
}

func (uc UseCase) AddEntry(sfID, deniedUserID entities.HexID) (DeniedEntry, error) {
	if uc.associatedUser.ID == deniedUserID {
		return DeniedEntry{}, apperr.Invalid(
			"denylist_self_entry",
			"you cannot add yourself to the denylist",
			nil,
		)
	}

	participant, err := uc.facProvider.participant.CheckParticipantInSecretFriend(
		sfID, uc.associatedUser.ID,
	)
	if err != nil {
		return DeniedEntry{}, apperr.Forbidden(
			"denylist_access_forbidden",
			"you are not a participant in this secret friend",
			err,
		)
	}

	if _, err = uc.facProvider.participant.CheckParticipantInSecretFriend(sfID, deniedUserID); err != nil {
		return DeniedEntry{}, apperr.Invalid(
			"denylist_target_not_participant",
			"target user is not a participant in this secret friend",
			err,
		)
	}

	// Validate capacity
	participantRef := ParticipantRef{
		ParticipantID:  participant.ID,
		UserID:         uc.associatedUser.ID,
		SecretFriendID: sfID,
	}

	currentList, err := uc.repo.GetDenyListByParticipant(participantRef)
	if err != nil {
		return DeniedEntry{}, apperr.From(
			"denylist_lookup_failed",
			"failed to load denylist",
			err,
		)
	}

	sf, err := uc.facProvider.secretFriend.GetSecretFriendByID(sfID)
	if err != nil {
		return DeniedEntry{}, apperr.From(
			"secret_friend_not_found",
			"secret friend not found",
			err,
		)
	}

	// Hard cap: denylist cannot exceed 50% of participants to guarantee a valid draw.
	// If MaxDenyListSize is also set, use the more restrictive limit.
	effectiveMax := int(sf.MaxDenyListSize)
	if participantCount := len(sf.Participants); participantCount > 1 {
		half := participantCount / 2
		if effectiveMax == 0 || effectiveMax > half {
			effectiveMax = half
		}
	}
	if effectiveMax > 0 && len(currentList) >= effectiveMax {
		return DeniedEntry{}, apperr.Conflict(
			"denylist_capacity_reached",
			"denylist capacity reached",
			nil,
		)
	}

	deniedUser, err := uc.repo.AddDenyListEntry(participantRef, deniedUserID)
	if err != nil {
		return DeniedEntry{}, apperr.From(
			"denylist_add_failed",
			"failed to add denylist entry",
			err,
		)
	}

	return deniedUser, nil
}

func (uc UseCase) RemoveEntry(sfID, deniedUserID entities.HexID) error {
	participant, err := uc.facProvider.participant.CheckParticipantInSecretFriend(
		sfID, uc.associatedUser.ID,
	)
	if err != nil {
		return apperr.Forbidden(
			"denylist_access_forbidden",
			"you are not a participant in this secret friend",
			err,
		)
	}

	if err = uc.repo.RemoveDenyListEntry(
		ParticipantRef{
			ParticipantID:  participant.ID,
			UserID:         uc.associatedUser.ID,
			SecretFriendID: sfID,
		},
		deniedUserID,
	); err != nil {
		return apperr.From(
			"denylist_remove_failed",
			"failed to remove denylist entry",
			err,
		)
	}

	return nil
}
