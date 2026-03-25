package changepassword

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/session"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type UseCase struct {
	passwordRepository Repository
	sessionRepository  session.Repository
}

func New(passwordRepository Repository, sessionRepository session.Repository) UseCase {
	return UseCase{
		passwordRepository: passwordRepository,
		sessionRepository:  sessionRepository,
	}
}

func (uc UseCase) Execute(
	authToken string,
	currentPassword string,
	newPassword string,
) error {
	user, err := session.FindAndCheckUser(uc.sessionRepository, authToken, currentPassword)
	if err != nil {
		return err
	}

	encryptedPassword, err := entities.UserBasic{Password: newPassword}.EncryptPassword()
	if err != nil {
		return autherrs.NewErrPasswordEncryption(err)
	}
	if err = uc.passwordRepository.UpdatePassword(user.ID, string(encryptedPassword)); err != nil {
		return autherrs.NewErrUpdatePassword(err)
	}

	return nil
}
