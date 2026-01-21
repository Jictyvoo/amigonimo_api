package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// ParticipantConverter is the converter for the entities.Participant type.
//
// goverter:converter
// goverter:output:file @cwd/participant_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend TimeFromNullTime
// goverter:extend CopyTime
// goverter:extend StringFromNullString
type ParticipantConverter interface {
	dbParticipantTimestampToEntity(p dbgen.Participant) entities.Timestamp

	// goverter:map UserID ID | HexIDFromBytes
	// goverter:ignore UserBasic
	// goverter:ignore FullName
	// goverter:ignore VerifiedAt
	// goverter:ignore RememberToken
	dbParticipantToRelatedUser(p dbgen.Participant) entities.User

	// goverter:map . Timestamp
	// goverter:map . RelatedUser
	// goverter:ignore DenyList
	// goverter:ignore Wishlist
	ToEntityParticipant(p dbgen.Participant) entities.Participant

	// goverter:ignore Password
	dbListParticipantsBySecretFriendRowToBasicUser(
		p dbgen.ListParticipantsBySecretFriendRow,
	) entities.UserBasic

	// goverter:map Fullname FullName
	// goverter:map UserID ID
	// goverter:map . UserBasic
	// goverter:ignore VerifiedAt
	// goverter:ignore RememberToken
	dbListParticipantsBySecretFriendRowToRelatedUser(
		p dbgen.ListParticipantsBySecretFriendRow,
	) entities.User

	// goverter:autoMap Participant
	// goverter:map Participant Timestamp
	// goverter:map . RelatedUser
	// goverter:ignore DenyList
	// goverter:ignore Wishlist
	MapParticipantRow(p dbgen.ListParticipantsBySecretFriendRow) entities.Participant
}
