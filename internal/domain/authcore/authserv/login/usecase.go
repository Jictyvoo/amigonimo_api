package login

import (
	"errors"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

const refreshTokenDuration = 24 * time.Hour * 60

type UseCase struct {
	userRepository  UserRepository
	tokenRepository TokenRepository
}

func New(userRepository UserRepository, tokenRepository TokenRepository) UseCase {
	return UseCase{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

func (uc UseCase) Execute(
	formUser entities.UserBasic,
) (authToken entities.AuthenticationToken, err error) {
	var user entities.User
	switch {
	case len(formUser.Email) > 0, len(formUser.Username) > 0:
		user, err = uc.userRepository.GetUserByEmailOrUsername(formUser.Email, formUser.Username)
	default:
		return authToken, autherrs.ErrInvalidCredentials
	}
	if err != nil && !errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return authToken, autherrs.NewErrLogin(err)
	}
	if user.ID.IsEmpty() {
		return authToken, autherrs.ErrInvalidCredentials
	}

	if ok, compareErr := user.ComparePassword(formUser.Password); !ok || compareErr != nil {
		return authToken, autherrs.ErrInvalidCredentials
	}

	if authToken, err = uc.tokenRepository.GetAuthenticationToken(user.ID); err != nil &&
		!errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return authToken, autherrs.NewErrTokenLookup(err)
	}

	if err = authToken.Regenerate(refreshTokenDuration); err != nil {
		return authToken, autherrs.NewErrTokenRegenerate(err)
	}

	if err = uc.tokenRepository.UpsertAuthToken(&authToken); err != nil {
		return authToken, autherrs.NewErrUpdateAuthToken(err)
	}

	authToken.User = user

	return authToken, nil
}
