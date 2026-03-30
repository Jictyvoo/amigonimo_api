package fixturesets

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

type User struct {
	User     *genmodels.User
	Profile  *genmodels.UserProfile
	Password string
}

func NewUser(email, password, fullname string) *User {
	builder := fixtures.NewUser().
		WithEmail(email).
		WithPassword(password)
	if fullname != "" {
		builder = builder.WithFullname(fullname)
	}

	return &User{
		User:     builder.Build(),
		Profile:  builder.BuildProfile(),
		Password: password,
	}
}

// WithRecoveryCode sets a recovery code and its expiry on the underlying user record.
func (u *User) WithRecoveryCode(code string, expiresAt time.Time) *User {
	u.User.RecoveryCode.Valid = true
	u.User.RecoveryCode.String = code
	u.User.RecoveryCodeExpiresAt.Valid = true
	u.User.RecoveryCodeExpiresAt.Time = expiresAt
	return u
}

// WithVerified marks the user as email-verified.
func (u *User) WithVerified(at time.Time) *User {
	u.User.VerifiedAt.Valid = true
	u.User.VerifiedAt.Time = at
	return u
}

func (u *User) ID() entities.HexID {
	return mustHexIDFromBytes(u.User.ID)
}
