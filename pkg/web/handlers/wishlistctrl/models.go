package wishlistctrl

import "time"

// WishlistItemResponse represents a wishlist item.
type WishlistItemResponse struct {
	ItemID   string    `json:"itemId"`
	Label    string    `json:"label"`
	Comments string    `json:"comments,omitempty"`
	AddedAt  time.Time `json:"addedAt,omitzero"`
}

// CreateWishlistItemRequest represents the request to create a wishlist item.
type CreateWishlistItemRequest struct {
	Label    string `json:"label"              validate:"required"`
	Comments string `json:"comments,omitempty"`
}
