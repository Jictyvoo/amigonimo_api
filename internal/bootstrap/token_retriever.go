package bootstrap

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/middlewares/jwtware"
)

func NewUserFromContext(ctx context.Context) (user entities.User, err error) {
	var claims jwt.MapClaims
	if claims, err = jwtware.ClaimsFromContext[jwt.MapClaims](ctx); err != nil {
		return user, err
	}

	// Parse user fields from claims
	{
		user.Username, _ = claims["username"].(string)
		verifiedAtUnix, _ := claims["verifiedAt"].(int)
		if verifiedAtUnix > 0 {
			user.VerifiedAt = time.UnixMilli(int64(verifiedAtUnix))
		}

		userIDStr, _ := claims["userID"].(string)
		userID, uuidErr := uuid.Parse(userIDStr)
		if uuidErr != nil {
			return user, uuidErr
		}

		user.ID = entities.HexID(userID)
	}

	return user, nil
}
