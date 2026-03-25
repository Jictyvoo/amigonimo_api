package requestrecovery

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
)

const resetCodeExpiration = 30 * time.Minute

type UseCase struct {
	userRepository UserRepository
	mailerService  Mailer
}

func New(userRepository UserRepository, mailer Mailer) UseCase {
	return UseCase{
		userRepository: userRepository,
		mailerService:  mailer,
	}
}

func (uc UseCase) Execute(userEmail string) error {
	user, err := uc.userRepository.GetUserByEmail(userEmail)
	if err != nil || user.ID.IsEmpty() {
		return autherrs.ErrUserEmailNotFound
	}

	recoveryCode, err := authcore.GenerateRecoveryCode(userEmail)
	if err != nil {
		return autherrs.NewErrGenRecoveryCode(err)
	}
	if err = uc.userRepository.SetRecoveryCode(user.ID, recoveryCode, time.Now().Add(resetCodeExpiration)); err != nil {
		return autherrs.NewErrGenRecoveryCode(err)
	}

	uc.mailerService.SendPasswordRecoveryEmail(user.Email, recoveryCode)

	return nil
}
