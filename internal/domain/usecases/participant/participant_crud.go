package participant

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) ConfirmParticipation(sfID entities.HexID) (entities.Participant, error) {
	// Validate that the Secret Friend group exists
	_, err := uc.secretFriendFacade.GetSecretFriendByID(sfID)
	if err != nil {
		return entities.Participant{}, fmt.Errorf("invalid secret friend group: %w", err)
	}

	return uc.repo.AddParticipant(sfID, uc.associatedUser.ID)
}

func (uc *UseCase) ListParticipants(sfID entities.HexID) ([]entities.Participant, error) {
	_, err := uc.repo.GetParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		isOwner, ownerErr := uc.secretFriendFacade.CheckUserIsOwner(sfID)
		if ownerErr != nil || !isOwner {
			return nil, fmt.Errorf("unauthorized: you are not a participant or owner of this group")
		}
	}
	return uc.repo.ListParticipants(sfID)
}

func (uc *UseCase) CheckParticipantExists(
	sfID, userID entities.HexID,
) (entities.Participant, error) {
	return uc.repo.GetParticipant(sfID, userID)
}

func (uc *UseCase) RemoveParticipant(sfID entities.HexID) error {
	_, err := uc.repo.GetParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		return fmt.Errorf("participant not found: %w", err)
	}
	return uc.repo.RemoveParticipant(sfID, uc.associatedUser.ID)
}
