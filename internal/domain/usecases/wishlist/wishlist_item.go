package wishlist

import (
	"strings"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// WishlistItem is a domain DTO owned by the wishlist package.
// It is stored per participant and returned by wishlist queries.
type WishlistItem struct {
	entities.Timestamp

	ID       entities.HexID
	Label    string
	Comments string
}

func (wi *WishlistItem) Normalize() {
	wi.Comments = strings.TrimSpace(wi.Comments)
	wi.Label = strings.TrimSpace(wi.Label)
	wi.Timestamp.Normalize()
}
