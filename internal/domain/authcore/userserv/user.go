package userserv

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type UserEditionService struct {
	userEditRepository UserEditionRepository
	userRepository     authserv.UserAuthRepository
	mailerService      authserv.MailerService
}

func NewUserEditService(
	userRepository authserv.UserAuthRepository,
	userEditRepository UserEditionRepository,
	mailer authserv.MailerService,
) UserEditionService {
	return UserEditionService{
		userRepository:     userRepository,
		userEditRepository: userEditRepository,
		mailerService:      mailer,
	}
}

func (serv UserEditionService) findAndCheckUser(
	authToken string, password string,
) (user entities.User, err error) {
	user, err = serv.userEditRepository.GetUserByAuthCode(authToken)
	if err != nil || user.ID.IsEmpty() {
		return entities.User{}, autherrs.ErrUserNotFound
	}

	// Check if current password matches stored password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return entities.User{}, autherrs.ErrWrongPassword
	}
	return user, nil
}

func (serv UserEditionService) ChangePassword(
	authToken string, currentPassword, newPassword string,
) error {
	user, err := serv.findAndCheckUser(authToken, currentPassword)
	if err != nil {
		return err
	}

	// Encrypt new password and save it on the database
	var encryptedPassword []byte
	if encryptedPassword, err = entities.UserBasic.EncryptPassword(entities.UserBasic{Password: newPassword}); err != nil {
		return autherrs.ErrPasswordEncryption
	}
	if err = serv.userRepository.UpdatePassword(user.ID, string(encryptedPassword)); err != nil {
		return autherrs.ErrUpdatePassword
	}
	return nil
}

func (serv UserEditionService) ChangeEmail(authToken string, emailForm entities.UserBasic) error {
	user, err := serv.findAndCheckUser(authToken, emailForm.Password)
	if err != nil {
		return err
	}
	if user.Email == emailForm.Email {
		return autherrs.ErrEmailInUse
	}

	existentUser, _ := serv.userRepository.GetUserByEmail(emailForm.Email)
	if !existentUser.ID.IsEmpty() {
		return autherrs.ErrEmailInUse
	}

	// Here the activation email will be generated
	verificationToken := authcore.GenerateActivationToken(
		user.Username + ":" + emailForm.Email,
	)

	// Change email and verification code
	if err = serv.userEditRepository.ChangeEmail(user.ID, emailForm.Email); err == nil {
		err = serv.userEditRepository.SetNewVerificationCode(user.ID, verificationToken)
	}

	// Send activation code in an email
	serv.mailerService.SendActivationEmail(emailForm.Email, verificationToken)
	return err
}

func (serv UserEditionService) ChangeUsername(
	authToken string, usernameForm entities.UserBasic,
) error {
	user, err := serv.findAndCheckUser(authToken, usernameForm.Password)
	if err != nil {
		return err
	}

	// Check if the username is not currently used
	if userWithUsername, _ := serv.userRepository.GetUserByUsername(usernameForm.Username); !userWithUsername.ID.IsEmpty() {
		return autherrs.ErrUsernameInUse
	}

	if err = serv.userEditRepository.UpdateUsername(user.ID, usernameForm.Username); err != nil {
		return autherrs.ErrUpdateUsername
	}
	return nil
}
