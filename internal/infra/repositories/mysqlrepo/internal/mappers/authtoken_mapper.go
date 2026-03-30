package mappers

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

// AuthTokenConverter is the converter for the entities.AuthenticationToken type.
//
// goverter:converter
// goverter:output:file @cwd/authtoken_mapper.gen.go
// goverter:output:format function
// goverter:extend HexIDFromBytes
// goverter:extend TimeFromNullTime
// goverter:extend StringFromNullString
// goverter:extend UUIDFromNullString
// goverter:extend CopyTime
type AuthTokenConverter interface {
	// goverter:map . BasicAuthToken
	// goverter:map . User
	ToEntityAuthenticationToken(source dbgen.AuthToken) entities.AuthenticationToken

	// Helper for the embedded BasicAuthToken
	// goverter:map Token AuthToken
	convertBasicAuthToken(source dbgen.AuthToken) entities.BasicAuthToken

	// goverter:map UserID ID
	// goverter:ignore VerifiedAt
	// goverter:ignore RememberToken
	// goverter:ignore UserBasic
	convertTokenUserToEntityUser(source dbgen.AuthToken) entities.User
}
