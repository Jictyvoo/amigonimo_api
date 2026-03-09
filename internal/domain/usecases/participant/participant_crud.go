package participant

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) ConfirmParticipation(sfID entities.HexID) (entities.Participant, error) {
	// TODO: Before performing the confirmation, it should check if the user was invited
	return uc.repo.AddParticipant(sfID, uc.associatedUser.ID)
}

func (uc *UseCase) ListParticipants(sfID entities.HexID) ([]entities.Participant, error) {
	return uc.repo.ListParticipants(sfID)
}

func (uc *UseCase) RemoveParticipant(sfID entities.HexID) error {
	_, err := uc.repo.GetParticipant(sfID, uc.associatedUser.ID)
	if err != nil {
		return fmt.Errorf("participant not found: %w", err)
	}
	return uc.repo.RemoveParticipant(sfID, uc.associatedUser.ID)
}
