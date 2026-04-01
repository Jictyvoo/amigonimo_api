package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// WishlistConverter maps DB rows to the wishlist.WishlistItem DTO.
//
// goverter:converter
// goverter:output:file @cwd/wishlist_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend CopyTime
// goverter:extend TimeFromNullTime
// goverter:extend StringFromNullString
type WishlistConverter interface {
	dbWishlistTimestampToEntity(p dbgen.WishlistItem) entities.Timestamp

	// goverter:map . Timestamp
	ToEntityWishlistItem(wi dbgen.WishlistItem) wishlist.WishlistItem
}
