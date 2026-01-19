package wishlistctrl

import (
	"context"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
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
		return nil, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	var loadedList []entities.WishlistItem
	if loadedList, err = uc.GetWishlist(sfID); err != nil {
		return nil, err
	}

	result := make([]WishlistItemResponse, len(loadedList))
	for index, item := range loadedList {
		result[index] = parseWishItem(item)
	}

	return result, nil
}

// CreateWishlistItem handles POST /wishlist.
func (h *Controller) CreateWishlistItem(
	c fuego.ContextWithBody[WishlistItemRequest],
) (WishlistItemResponse, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return WishlistItemResponse{}, err
	}

	body, decodeErr := c.Body()
	if decodeErr != nil {
		return WishlistItemResponse{}, decodeErr
	}

	uc, facErr := h.useCaseFactory(c.Context())
	if facErr != nil {
		return WishlistItemResponse{}, facErr
	}

	addedItem, err := uc.AddItem(sfID, body.Label, body.Comments)
	if err != nil {
		return WishlistItemResponse{}, err
	}

	return parseWishItem(addedItem), nil
}

// DeleteWishlistItem handles DELETE /wishlist/{itemId}.
func (h *Controller) DeleteWishlistItem(
	c fuego.ContextNoBody,
) (any, error) {
	sfID, err := h.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	itemIDStr := c.PathParam("itemId")
	itemID, err := entities.ParseHexID(itemIDStr)
	if err != nil {
		return nil, err
	}

	uc, err := h.useCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	return nil, uc.DeleteItem(sfID, itemID)
}
