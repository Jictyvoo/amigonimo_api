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

	summaries, err := uc.ListSummaries(sfID)
	if err != nil {
		return nil, h.HandleError(err)
	}

	resp := make([]ParticipantResponse, len(summaries))
	for i, s := range summaries {
		resp[i] = ParticipantResponse{
			ParticipantID: s.ID.String(),
			UserID:        s.RelatedUser.ID.String(),
			Fullname:      s.FullName,
		}
	}

	return resp, nil
}
