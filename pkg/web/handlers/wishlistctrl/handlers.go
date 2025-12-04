package wishlistctrl

import (
	"github.com/go-fuego/fuego"
)

type WishlistHandlers struct {
	// TODO: Add service dependencies
}

func NewWishlistHandlers() *WishlistHandlers {
	return &WishlistHandlers{}
}

// GetWishlist handles GET /secret-friends/{id}/wishlist
func (h *WishlistHandlers) GetWishlist(
	c fuego.ContextNoBody,
) ([]WishlistItemResponse, error) {
	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	return nil, nil
}

// CreateWishlistItem handles POST /secret-friends/{id}/wishlist
func (h *WishlistHandlers) CreateWishlistItem(
	c fuego.ContextWithBody[CreateWishlistItemRequest],
) (*WishlistItemResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	_ = req
	return nil, nil
}

// DeleteWishlistItem handles DELETE /secret-friends/{id}/wishlist/{itemId}
func (h *WishlistHandlers) DeleteWishlistItem(
	c fuego.ContextNoBody,
) (any, error) {
	// TODO: Extract secretFriendId and itemId from path
	// TODO: Implement service call
	return nil, nil
}
