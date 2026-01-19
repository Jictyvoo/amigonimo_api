package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// WishlistConverter is the converter for the entities.WishlistItem type.
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
	ToEntityWishlistItem(wi dbgen.WishlistItem) entities.WishlistItem
}
