package denylistctrl

import (
	"context"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type (
	UseCaseFactory[T any] func(ctx context.Context) (T, error)
	Controller            struct {
		web.DefaultController

		useCaseFactory UseCaseFactory[denylist.UseCase]
	}
)

func NewController(useCaseFac UseCaseFactory[denylist.UseCase]) *Controller {
	return &Controller{useCaseFactory: useCaseFac}
}

type DenyListEntryRequest struct {
	DeniedUserID string `json:"deniedUserId"`
}

// GetDenyList handles GET /denylist.
func (h *Controller) GetDenyList(
	c fuego.ContextNoBody,
) ([]entities.DeniedUser, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	return uc.GetDenyList(sfID)
}

// AddDenyListEntry handles POST /denylist.
func (h *Controller) AddDenyListEntry(
	c fuego.ContextWithBody[DenyListEntryRequest],
) (entities.DeniedUser, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return entities.DeniedUser{}, err
	}

	body, err := c.Body()
	if err != nil {
		return entities.DeniedUser{}, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return entities.DeniedUser{}, err
	}

	deniedUserID, err := entities.ParseHexID(body.DeniedUserID)
	if err != nil {
		return entities.DeniedUser{}, err
	}

	return uc.AddEntry(sfID, deniedUserID)
}

// RemoveDenyListEntry handles DELETE /denylist/{deniedUserId}.
func (h *Controller) RemoveDenyListEntry(
	c fuego.ContextNoBody,
) (any, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	deniedUserIDStr := c.PathParam("deniedUserId")
	deniedUserID, err := entities.ParseHexID(deniedUserIDStr)
	if err != nil {
		return nil, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	return nil, uc.RemoveEntry(sfID, deniedUserID)
}
