package facades

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// Ensure ParticipantFacade implements required interfaces.
var (
	_ wishlist.ParticipantFacade = (*ParticipantFacade)(nil)
	_ denylist.ParticipantFacade = (*ParticipantFacade)(nil)
)

type ParticipantFacade struct {
	ports.BaseFacade
	participantUC participant.UseCase
}

func NewParticipantFacade(
	participantUC participant.UseCase,
) *ParticipantFacade {
	return &ParticipantFacade{
		participantUC: participantUC,
	}
}

// CheckParticipantInSecretFriend implements wishlist.ParticipantFacade and denylist.ParticipantFacade.
func (f *ParticipantFacade) CheckParticipantInSecretFriend(
	sfID, userID entities.HexID,
) (entities.Participant, error) {
	p, err := f.participantUC.CheckParticipantExists(sfID, userID)
	if err != nil {
		return entities.Participant{}, fmt.Errorf("user is not a participant: %w", err)
	}

	return p, nil
}
