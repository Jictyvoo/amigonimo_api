package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BasicAuthToken struct {
	AuthToken    string
	ExpiresAt    time.Time
	RefreshToken uuid.NullUUID
}

type AuthenticationToken struct {
	ID   HexID
	User User
	BasicAuthToken
}

type VerifyToken struct {
	User  User
	Token string
}

func (bat *BasicAuthToken) Regenerate(tokenDuration time.Duration) error {
	{
		u, err := uuid.NewRandom()
		if err != nil {
			u = uuid.New()
		}
		bat.AuthToken = u.String()
	}

	newRefreshToken, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("could not generate refresh token: %w", err)
	}

	// create a refresh token using uuid
	bat.RefreshToken = uuid.NullUUID{
		UUID:  newRefreshToken,
		Valid: len(newRefreshToken.String()) > 0,
	}

	// update expiration time
	bat.ExpiresAt = time.Now().Add(tokenDuration)
	return nil
}
