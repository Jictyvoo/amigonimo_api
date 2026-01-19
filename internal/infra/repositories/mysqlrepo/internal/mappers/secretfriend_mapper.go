package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// SecretFriendConverter is the converter for the entities.SecretFriend type.
//
// goverter:converter
// goverter:output:file @cwd/secretfriend_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend TimeFromNullTime
// goverter:extend StringFromNullString
// goverter:extend CopyTime
type SecretFriendConverter interface {
	dbSecretFriendTimestampToEntity(p dbgen.SecretFriend) entities.Timestamp

	// goverter:map Status | StatusToEntity
	// goverter:map . Timestamp
	// goverter:ignore Participants
	// goverter:ignore DrawResult
	ToEntitySecretFriend(sf dbgen.SecretFriend) entities.SecretFriend
}

func StatusToEntity(status string) entities.SecretFriendStatus {
	return entities.SecretFriendStatus(status)
}

func DatetimeToEntity(dt interface{}) *interface{} {
	// This is a placeholder for goverter logic if needed for pointers
	return nil
}
