package usecases

import (
	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type DrawUseCase struct{}

func NewDrawUseCase() *DrawUseCase {
	return &DrawUseCase{}
}

// GetDrawResult retrieves the draw result for a participant.
func (uc *DrawUseCase) GetDrawResult(
	userID, secretFriendID uuid.UUID,
) (*entities.DrawResult, error) {
	return nil, nil
}
