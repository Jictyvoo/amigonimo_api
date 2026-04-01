package resetpassword

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/checkrecovery"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

type UseCase struct {
	userRepository Repository
}

func New(userRepository Repository) UseCase {
	return UseCase{userRepository: userRepository}
}

func (uc UseCase) Execute(resetUser authvalues.UserBasic, recoveryCode string) error {
	checkRecovery := checkrecovery.New(uc.userRepository)
	userID, err := checkRecovery.Execute(resetUser.Email, recoveryCode)
	if err != nil {
		return err
	}

	encryptedPassword, encryptErr := resetUser.EncryptPassword()
	if encryptErr != nil {
		return autherrs.NewErrPasswordEncryption(encryptErr)
	}
	if err = uc.userRepository.UpdatePassword(userID, string(encryptedPassword)); err != nil {
		return autherrs.NewErrUpdatePassword(err)
	}

	_ = uc.userRepository.SetRecoveryCode(userID, "", time.Time{})

	return nil
}
