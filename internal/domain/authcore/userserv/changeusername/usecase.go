package changeusername

import (
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/session"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

type UseCase struct {
	userRepository     UserRepository
	userEditRepository Repository
}

func New(userRepository UserRepository, userEditRepository Repository) UseCase {
	return UseCase{
		userRepository:     userRepository,
		userEditRepository: userEditRepository,
	}
}

func (uc UseCase) Execute(
	authToken string,
	usernameForm authvalues.UserBasic,
) error {
	user, err := session.FindAndCheckUser(
		uc.userEditRepository,
		authToken,
		usernameForm.Password,
	)
	if err != nil {
		return err
	}

	userWithUsername, err := uc.userRepository.GetUserByUsername(usernameForm.Username)
	if err != nil && !errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return autherrs.NewErrChangeUsernameLookup(err)
	}
	if !userWithUsername.ID.IsEmpty() {
		return autherrs.ErrUsernameInUse
	}

	if err = uc.userEditRepository.UpdateUsername(user.ID, usernameForm.Username); err != nil {
		return autherrs.NewErrUpdateUsername(err)
	}

	return nil
}
