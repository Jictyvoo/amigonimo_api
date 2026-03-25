package lookuprecovery

import (
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

type UseCase struct {
	userRepository Repository
}

func New(userRepository Repository) UseCase {
	return UseCase{userRepository: userRepository}
}

func (uc UseCase) Execute(username string) (string, error) {
	user, err := uc.userRepository.GetUserByUsername(username)
	if err != nil && !errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return "", autherrs.NewErrRecoveryLookup(err)
	}
	if user.ID.IsEmpty() {
		return "", autherrs.ErrUserNotFound
	}

	return user.ObfuscateEmail(), nil
}
