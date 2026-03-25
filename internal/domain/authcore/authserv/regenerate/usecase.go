package regenerate

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

const refreshTokenDuration = 24 * time.Hour * 60

type UseCase struct {
	tokenRepository Repository
}

func New(tokenRepository Repository) UseCase {
	return UseCase{tokenRepository: tokenRepository}
}

func (uc UseCase) Execute(refreshToken string) (entities.AuthenticationToken, error) {
	authentication, err := uc.tokenRepository.CheckAuthenticationByRefreshToken(refreshToken)
	if err != nil || authentication.User.ID.IsEmpty() ||
		time.Now().After(authentication.ExpiresAt) {
		return entities.AuthenticationToken{}, autherrs.ErrInvalidAuthToken
	}

	if err = authentication.Regenerate(refreshTokenDuration); err != nil {
		return entities.AuthenticationToken{}, autherrs.NewErrTokenRegenerate(err)
	}

	if err = uc.tokenRepository.UpsertAuthToken(&authentication); err != nil {
		return entities.AuthenticationToken{}, autherrs.NewErrUpdateAuthToken(err)
	}

	return authentication, nil
}
