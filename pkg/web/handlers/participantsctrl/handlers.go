package participantsctrl

import (
	"github.com/go-fuego/fuego"
)

type ParticipantsHandlers struct {
	// TODO: Add service dependencies
}

func NewParticipantsHandlers() *ParticipantsHandlers {
	return &ParticipantsHandlers{}
}

// ConfirmParticipation handles POST /secret-friends/{id}/participants.
func (h *ParticipantsHandlers) ConfirmParticipation(
	c fuego.ContextWithBody[ConfirmParticipationRequest],
) (*ConfirmParticipationResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	_ = req
	return nil, nil
}

// ListParticipants handles GET /secret-friends/{id}/participants.
func (h *ParticipantsHandlers) ListParticipants(
	c fuego.ContextNoBody,
) ([]ParticipantResponse, error) {
	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	return nil, nil
}
