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
	// goverter:ignore Password
	dbDenyRowToBasicDeniedUserData(row dbgen.GetDenyListByParticipantRow) entities.UserBasic

	// goverter:map Fullname FullName
	// goverter:map . UserBasic
	// goverter:map Denylist.DeniedUserID ID
	// goverter:ignore VerifiedAt
	// goverter:ignore RememberToken
	dbDenyRowToDeniedUser(row dbgen.GetDenyListByParticipantRow) entities.User

	// goverter:map Denylist.ParticipantID ID
	// goverter:map Denylist.CreatedAt JoinedAt
	// goverter:map . RelatedUser
	// goverter:ignore SecretFriendID
	// goverter:ignore Timestamp
	// goverter:ignore DenyList
	// goverter:ignore Wishlist
	dbDenyRowToDeniedParticipant(row dbgen.GetDenyListByParticipantRow) entities.Participant

	// goverter:map Denylist.ID ID
	// goverter:map Denylist Timestamp
	// goverter:map . InnerParticipant
	ToEntityDeniedUser(row dbgen.GetDenyListByParticipantRow) entities.DeniedUser
}
