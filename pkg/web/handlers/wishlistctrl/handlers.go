package wishlistctrl

import (
	"context"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type (
	UseCaseFactory[T any] func(ctx context.Context) (T, error)
	Controller            struct {
		web.DefaultController

		useCaseFactory UseCaseFactory[wishlist.UseCase]
	}
)

func NewController(useCaseFac UseCaseFactory[wishlist.UseCase]) *Controller {
	return &Controller{useCaseFactory: useCaseFac}
}

type WishlistItemRequest struct {
	Label    string `json:"label"`
	Comments string `json:"comments"`
}

// GetWishlist handles GET /wishlist.
func (h *Controller) GetWishlist(
	c fuego.ContextNoBody,
) ([]WishlistItemResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, h.HandleError(err)
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, h.HandleError(err)
	}

	var loadedList []wishlist.WishlistItem
	if loadedList, err = uc.GetWishlist(sfID); err != nil {
		return nil, h.HandleError(err)
	}

	result := make([]WishlistItemResponse, len(loadedList))
	for index, item := range loadedList {
		result[index] = parseWishItem(item)
	}

	return result, nil
}

// CreateWishlistItem handles POST /wishlist.
func (h *Controller) CreateWishlistItem(
	c fuego.Context[WishlistItemRequest, struct{}],
) (WishlistItemResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return WishlistItemResponse{}, h.HandleError(err)
	}

	body, decodeErr := c.Body()
	if decodeErr != nil {
		return WishlistItemResponse{}, h.HandleError(decodeErr)
	}

	uc, facErr := h.useCaseFactory(c.Context())
	if facErr != nil {
		return WishlistItemResponse{}, h.HandleError(facErr)
	}

	addedItem, err := uc.AddItem(sfID, body.Label, body.Comments)
	if err != nil {
		return WishlistItemResponse{}, h.HandleError(err)
	}

	return parseWishItem(addedItem), nil
}

// DeleteWishlistItem handles DELETE /wishlist/{itemId}.
func (h *Controller) DeleteWishlistItem(
	c fuego.ContextNoBody,
) (DeleteWishlistItemResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return DeleteWishlistItemResponse{}, h.HandleError(err)
	}

	itemIDStr := c.PathParam("itemId")
	itemID, err := h.ParseHexID(itemIDStr)
	if err != nil {
		return DeleteWishlistItemResponse{}, h.HandleError(err)
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return DeleteWishlistItemResponse{}, h.HandleError(err)
	}

	if err = uc.DeleteItem(sfID, itemID); err != nil {
		return DeleteWishlistItemResponse{}, h.HandleError(err)
	}

	return DeleteWishlistItemResponse{
		Success:   true,
		DeletedID: itemID.String(),
	}, nil
}
