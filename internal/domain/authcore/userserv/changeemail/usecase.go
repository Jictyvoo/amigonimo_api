package changeemail

import (
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/session"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

type UseCase struct {
	userRepository     UserRepository
	userEditRepository Repository
	mailerService      Mailer
}

func New(userRepository UserRepository, userEditRepository Repository, mailer Mailer) UseCase {
	return UseCase{
		userRepository:     userRepository,
		userEditRepository: userEditRepository,
		mailerService:      mailer,
	}
}

func (uc UseCase) Execute(authToken string, emailForm authvalues.UserBasic) error {
	user, err := session.FindAndCheckUser(
		uc.userEditRepository,
		authToken,
		emailForm.Password,
	)
	if err != nil {
		return err
	}
	if user.Email == emailForm.Email {
		return autherrs.ErrEmailInUse
	}

	existentUser, err := uc.userRepository.GetUserByEmail(emailForm.Email)
	if err != nil && !errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return autherrs.NewErrChangeEmailLookup(err)
	}
	if !existentUser.ID.IsEmpty() {
		return autherrs.ErrEmailInUse
	}

	verificationToken := authcore.GenerateActivationToken(user.Username + ":" + emailForm.Email)
	if err = uc.userEditRepository.ChangeEmail(user.ID, emailForm.Email); err != nil {
		return autherrs.NewErrUpdateEmail(err)
	}
	if err = uc.userEditRepository.SetNewVerificationCode(user.ID, verificationToken); err != nil {
		return autherrs.NewErrSetVerification(err)
	}

	uc.mailerService.SendActivationEmail(emailForm.Email, verificationToken)

	return nil
}
