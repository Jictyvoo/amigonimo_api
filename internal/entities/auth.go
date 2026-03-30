package entities

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

// UserBasic is a type alias for backward compatibility.
type UserBasic = authvalues.UserBasic

type User struct {
	authvalues.UserBasic

	ID            HexID
	VerifiedAt    time.Time
	RememberToken string
}
