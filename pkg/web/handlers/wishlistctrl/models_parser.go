package wishlistctrl

import "github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"

func parseWishItem(item wishlist.WishlistItem) WishlistItemResponse {
	return WishlistItemResponse{
		ItemID:   item.ID.String(),
		Label:    item.Label,
		Comments: item.Comments,
		AddedAt:  item.CreatedAt,
	}
}
