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
	dbPartipantTimestampToEntity(p dbgen.Participant) entities.Timestamp

	// goverter:map UserID ID | HexIDFromBytes
	// goverter:ignore UserBasic
	// goverter:ignore FullName
	// goverter:ignore VerifiedAt
	// goverter:ignore RememberToken
	dbPartipantToRelatedUser(p dbgen.Participant) entities.User

	// goverter:map . Timestamp
	// goverter:map . RelatedUser
	// goverter:ignore DenyList
	// goverter:ignore Wishlist
	ToEntityParticipant(p dbgen.Participant) entities.Participant
}

func MapParticipantRow(p dbgen.ListParticipantsBySecretFriendRow) entities.Participant {
	return entities.Participant{
		Timestamp: entities.Timestamp{
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		ID:             HexIDFromBytes(p.ID),
		SecretFriendID: HexIDFromBytes(p.SecretFriendID),
		JoinedAt:       p.JoinedAt.Time,
		RelatedUser: entities.User{
			ID:       HexIDFromBytes(p.UserID),
			FullName: p.Fullname,
			UserBasic: entities.UserBasic{
				Email:    p.Email,
				Username: p.Username,
			},
		},
	}
}
