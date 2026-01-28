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

// GetDenyList handles GET /denylist.
func (h *Controller) GetDenyList(
	c fuego.ContextNoBody,
) ([]DeniedUserResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	deniedUsers, err := uc.GetDenyList(sfID)
	if err != nil {
		return nil, err
	}

	result := make([]DeniedUserResponse, len(deniedUsers))
	for i, user := range deniedUsers {
		result[i] = parseDeniedUser(user)
	}

	return result, nil
}

// AddDenyListEntry handles POST /denylist.
func (h *Controller) AddDenyListEntry(
	c fuego.Context[AddDenyListRequest, struct{}],
) (DeniedUserResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return DeniedUserResponse{}, err
	}

	body, err := c.Body()
	if err != nil {
		return DeniedUserResponse{}, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return DeniedUserResponse{}, err
	}

	deniedUserID, err := entities.ParseHexID(body.TargetUserID)
	if err != nil {
		return DeniedUserResponse{}, err
	}

	deniedUser, err := uc.AddEntry(sfID, deniedUserID)
	if err != nil {
		return DeniedUserResponse{}, err
	}

	return parseDeniedUser(deniedUser), nil
}

// RemoveDenyListEntry handles DELETE /denylist/{deniedUserId}.
func (h *Controller) RemoveDenyListEntry(
	c fuego.ContextNoBody,
) (RemoveDenyListEntryResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return RemoveDenyListEntryResponse{}, err
	}

	deniedUserIDStr := c.PathParam("deniedUserId")
	deniedUserID, err := entities.ParseHexID(deniedUserIDStr)
	if err != nil {
		return RemoveDenyListEntryResponse{}, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return RemoveDenyListEntryResponse{}, err
	}

	if err := uc.RemoveEntry(sfID, deniedUserID); err != nil {
		return RemoveDenyListEntryResponse{}, err
	}

	return RemoveDenyListEntryResponse{
		Success:   true,
		DeletedID: deniedUserID.String(),
	}, nil
}
