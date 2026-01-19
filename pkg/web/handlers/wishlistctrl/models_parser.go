package wishlistctrl

import "github.com/jictyvoo/amigonimo_api/internal/entities"

func parseWishItem(item entities.WishlistItem) WishlistItemResponse {
	return WishlistItemResponse{
		ItemID:   item.ID.String(),
		Label:    item.Label,
		Comments: item.Comments,
		AddedAt:  item.CreatedAt,
	}
}
