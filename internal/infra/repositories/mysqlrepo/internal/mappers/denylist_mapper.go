package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// DenylistConverter is the converter for the entities.DeniedUser type.
//
// goverter:converter
// goverter:output:file @cwd/denylist_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend CopyTime
type DenylistConverter interface {
	// goverter:map . RelatedUser
	// goverter:map . Timestamp
	// goverter:map CreatedAt JoinedAt
	// goverter:ignore SecretFriendID
	// goverter:ignore DenyList
	// goverter:ignore Wishlist
	dbDenyRowToDeniedParticipant(row dbgen.GetDenyListByParticipantRow) entities.Participant

	// goverter:map DeniedUserID ID
	// goverter:map Fullname FullName
	// goverter:map . UserBasic
	// goverter:ignore VerifiedAt
	// goverter:ignore RememberToken
	dbDenyRowToDeniedUser(row dbgen.GetDenyListByParticipantRow) entities.User

	// goverter:ignore Password
	dbDenyRowToBasicDeniedUserData(row dbgen.GetDenyListByParticipantRow) entities.UserBasic

	// goverter:map . Timestamp
	// goverter:map . InnerParticipant
	// goverter:map DeniedUserID ID
	ToEntityDeniedUser(row dbgen.GetDenyListByParticipantRow) entities.DeniedUser
}
