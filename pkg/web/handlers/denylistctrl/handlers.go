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
		return nil, h.HandleError(err)
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, h.HandleError(err)
	}

	deniedUsers, err := uc.GetDenyList(sfID)
	if err != nil {
		return nil, h.HandleError(err)
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
		return DeniedUserResponse{}, h.HandleError(err)
	}

	body, decodeErr := c.Body()
	if decodeErr != nil {
		return DeniedUserResponse{}, h.HandleError(decodeErr)
	}

	var deniedUserID entities.HexID
	if deniedUserID, err = h.ParseHexID(body.TargetUserID); err != nil {
		return DeniedUserResponse{}, h.HandleError(err)
	}

	uc, ucErr := h.useCaseFactory(c.Context())
	if ucErr != nil {
		return DeniedUserResponse{}, h.HandleError(ucErr)
	}

	deniedUser, err := uc.AddEntry(sfID, deniedUserID)
	if err != nil {
		return DeniedUserResponse{}, h.HandleError(err)
	}

	return parseDeniedUser(deniedUser), nil
}

// RemoveDenyListEntry handles DELETE /denylist/{deniedUserId}.
func (h *Controller) RemoveDenyListEntry(
	c fuego.ContextNoBody,
) (RemoveDenyListEntryResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return RemoveDenyListEntryResponse{}, h.HandleError(err)
	}

	var deniedUserID entities.HexID
	if deniedUserID, err = h.ParseHexID(c.PathParam("deniedUserId")); err != nil {
		return RemoveDenyListEntryResponse{}, h.HandleError(err)
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return RemoveDenyListEntryResponse{}, h.HandleError(err)
	}

	if err = uc.RemoveEntry(sfID, deniedUserID); err != nil {
		return RemoveDenyListEntryResponse{}, h.HandleError(err)
	}

	return RemoveDenyListEntryResponse{
		Success:   true,
		DeletedID: deniedUserID.String(),
	}, nil
}
