package signup

import (
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

type UseCase struct {
	userRepository Repository
	mailerService  Mailer
}

func New(userRepository Repository, mailer Mailer) UseCase {
	return UseCase{
		userRepository: userRepository,
		mailerService:  mailer,
	}
}

func (uc UseCase) Execute(inputUser authvalues.UserBasic) (entities.User, error) {
	user, err := uc.userRepository.GetUserByEmailOrUsername(inputUser.Email, inputUser.Username)
	if err != nil && !errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return entities.User{}, autherrs.NewErrSignUpLookup(err)
	}
	if !user.ID.IsEmpty() {
		return entities.User{}, autherrs.ErrEmailOrUsernameUsed
	}

	encryptedPassword, err := inputUser.EncryptPassword()
	if err != nil {
		return entities.User{}, autherrs.NewErrPasswordEncryption(err)
	}

	newUser := entities.User{
		UserBasic: authvalues.UserBasic{
			Username: inputUser.Username,
			Email:    inputUser.Email,
			Password: string(encryptedPassword),
		},
		RememberToken: "",
	}
	if newUser.ID, err = entities.NewHexID(); err != nil {
		return entities.User{}, autherrs.NewErrUserCreation(err)
	}

	verificationToken := authcore.GenerateActivationToken(newUser.ID.String())
	if err = uc.userRepository.CreateUser(newUser, verificationToken); err != nil {
		return entities.User{}, autherrs.NewErrUserCreation(err)
	}

	uc.mailerService.SendActivationEmail(newUser.Email, verificationToken)

	return newUser, nil
}
