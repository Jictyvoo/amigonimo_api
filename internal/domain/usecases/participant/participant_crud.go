package participant

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) ConfirmParticipation(sfID entities.HexID) (entities.Participant, error) {
	// Validate that the Secret Friend group exists
	_, err := uc.secretFriendFacade.GetSecretFriendByID(sfID)
	if err != nil {
		return entities.Participant{}, apperr.From(
			"secret_friend_not_found",
			"secret friend not found",
			err,
		)
	}

	participant, err := uc.repo.AddParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		return entities.Participant{}, apperr.From(
			"participant_confirm_failed",
			"failed to confirm participation",
			err,
		)
	}

	return participant, nil
}

func (uc *UseCase) ListParticipants(sfID entities.HexID) ([]entities.Participant, error) {
	_, err := uc.repo.GetParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		isOwner, ownerErr := uc.secretFriendFacade.CheckUserIsOwner(sfID)
		if ownerErr != nil {
			return nil, ownerErr
		}
		if !isOwner {
			return nil, apperr.Forbidden(
				"participant_list_forbidden",
				"you are not allowed to view this participant list",
				nil,
			)
		}
	}

	participants, err := uc.repo.ListParticipants(sfID)
	if err != nil {
		return nil, apperr.From(
			"participant_list_failed",
			"failed to list participants",
			err,
		)
	}

	return participants, nil
}

func (uc *UseCase) MarkAsReady(sfID entities.HexID) error {
	currentParticipant, err := uc.repo.GetParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		return apperr.From(
			"participant_not_found",
			"participant not found",
			err,
		)
	}

	if currentParticipant.IsReady {
		return nil
	}

	if err = uc.repo.SetParticipantReady(sfID, uc.associatedUser.ID, true); err != nil {
		return apperr.From(
			"participant_ready_update_failed",
			"failed to update participant readiness",
			err,
		)
	}

	return nil
}

func (uc *UseCase) CheckParticipantExists(
	sfID, userID entities.HexID,
) (entities.Participant, error) {
	return uc.repo.GetParticipant(sfID, userID)
}

func (uc *UseCase) RemoveParticipant(sfID entities.HexID) error {
	_, err := uc.repo.GetParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		return apperr.From(
			"participant_not_found",
			"participant not found",
			err,
		)
	}

	if err = uc.repo.RemoveParticipant(sfID, uc.associatedUser.ID); err != nil {
		return apperr.From(
			"participant_remove_failed",
			"failed to remove participant",
			err,
		)
	}

	return nil
}
