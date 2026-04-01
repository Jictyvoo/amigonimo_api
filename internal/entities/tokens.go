package entities

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

type AuthenticationToken struct {
	authvalues.BasicAuthToken

	ID   HexID
	User User
}
