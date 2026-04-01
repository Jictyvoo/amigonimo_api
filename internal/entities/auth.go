package entities

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

type User struct {
	authvalues.UserBasic

	ID            HexID
	VerifiedAt    time.Time
	RememberToken string
}
