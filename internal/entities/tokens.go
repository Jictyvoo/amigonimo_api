package entities

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

// BasicAuthToken is a type alias for backward compatibility.
type BasicAuthToken = authvalues.BasicAuthToken

type AuthenticationToken struct {
	authvalues.BasicAuthToken

	ID   HexID
	User User
}
