package participantsctrl

import (
	"context"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type (
	UseCaseFactory[T any] func(ctx context.Context) (T, error)
	ParticipantsHandlers  struct {
		web.DefaultController

		useCaseFactory UseCaseFactory[participant.UseCase]
	}
)

func NewParticipantsHandlers(
	useCaseFac UseCaseFactory[participant.UseCase],
) *ParticipantsHandlers {
	return &ParticipantsHandlers{useCaseFactory: useCaseFac}
}

// ConfirmParticipation handles POST /secret-friends/{id}/participants.
func (h *ParticipantsHandlers) ConfirmParticipation(
	c fuego.Context[ConfirmParticipationRequest, struct{}],
) (*ConfirmParticipationResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, h.HandleError(err)
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, h.HandleError(err)
	}

	p, err := uc.ConfirmParticipation(sfID)
	if err != nil {
		return nil, h.HandleError(err)
	}

	return &ConfirmParticipationResponse{
		Success:       true,
		ParticipantID: p.ID.String(),
	}, nil
}

// ListParticipants handles GET /secret-friends/{id}/participants.
func (h *ParticipantsHandlers) ListParticipants(
	c fuego.ContextNoBody,
) ([]ParticipantResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, h.HandleError(err)
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, h.HandleError(err)
	}

	participants, err := uc.ListParticipants(sfID)
	if err != nil {
		return nil, h.HandleError(err)
	}

	resp := make([]ParticipantResponse, len(participants))
	for i, p := range participants {
		resp[i] = ParticipantResponse{
			ParticipantID: p.ID.String(),
			UserID:        p.RelatedUser.ID.String(),
			Fullname:      p.RelatedUser.FullName,
		}
	}

	return resp, nil
}
