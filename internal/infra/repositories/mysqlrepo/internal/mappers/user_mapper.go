package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// UserConverter is the converter for the entities.User type.
//
// goverter:converter
// goverter:output:file @cwd/user_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend TimeFromNullTime
// goverter:extend StringFromNullString
// goverter:extend UUIDFromNullString
type UserConverter interface {
	// Helper method to convert UserBasic fields - will be used to populate embedded UserBasic
	convertUserBasic(user dbgen.User) entities.UserBasic

	// goverter:map . UserBasic
	ToEntityUser(user dbgen.User) entities.User
}
