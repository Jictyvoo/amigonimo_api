package authserv

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

const (
	refreshTokenDuration = 24 * time.Hour * 60
	resetCodeExpiration  = 30 * time.Minute
)

type AuthService struct {
	userRepository  UserAuthRepository
	tokenRepository TokenRepository
	mailerService   MailerService
}

func NewAuthService(
	userRepository UserAuthRepository,
	tokenRepository TokenRepository,
	mailer MailerService,
) AuthService {
	return AuthService{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		mailerService:   mailer,
	}
}

func (serv AuthService) UserSignUp(inputUser entities.UserBasic) (entities.User, error) {
	if user, err := serv.userRepository.GetUserByEmailOrUsername(inputUser.Email, inputUser.Username); err == nil &&
		!user.ID.IsEmpty() {
		return entities.User{}, autherrs.ErrEmailOrUsernameUsed
	}
	encryptedPassword, err := inputUser.EncryptPassword()
	if err != nil {
		return entities.User{}, autherrs.ErrPasswordEncryption
	}
	newUser := entities.User{
		ID: entities.HexID(uuid.New()),
		UserBasic: entities.UserBasic{
			Username: inputUser.Username,
			Email:    inputUser.Email,
			Password: string(encryptedPassword),
		},
		RememberToken: "",
	}
	// Here the activation email will be generated
	verificationToken := authcore.GenerateActivationToken(newUser.Username + ":" + newUser.Email)
	if err = serv.userRepository.CreateUser(newUser, verificationToken); err != nil {
		return entities.User{}, autherrs.ErrUserCreation
	}
	serv.mailerService.SendActivationEmail(newUser.Email, verificationToken)
	return newUser, nil
}

func (serv AuthService) UserLogIn(
	formUser entities.UserBasic,
) (authToken entities.AuthenticationToken, err error) {
	var user entities.User
	switch {
	case len(formUser.Email) > 0, len(formUser.Username) > 0:
		user, err = serv.userRepository.GetUserByEmailOrUsername(formUser.Email, formUser.Username)
	default:
		return authToken, autherrs.ErrUserEmailNotFound
	}
	if err != nil && !errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return authToken, fmt.Errorf("internal server error: %w", err)
	}

	if user.ID.IsEmpty() {
		return authToken, autherrs.ErrUserEmailNotFound
	}
	if ok, compareErr := user.ComparePassword(formUser.Password); !ok || compareErr != nil {
		return authToken, autherrs.ErrWrongPassword
	}

	if authToken, err = serv.tokenRepository.GetAuthenticationToken(user.ID); err != nil &&
		!errors.Is(err, &dberrs.ErrDatabaseNotFound{}) {
		return authToken, fmt.Errorf("internal server error: %w", err)
	}

	// Update refresh token and expiration only (JWT will be generated in controller)
	if err = authToken.Regenerate(refreshTokenDuration); err != nil {
		return authToken, fmt.Errorf("failed to regenerate refresh token: %w", err)
	}

	// Save AuthenticationToken update
	if err = serv.tokenRepository.UpsertAuthToken(&authToken); err != nil {
		return authToken, autherrs.ErrUpdateAuthToken
	}

	authToken.User = user // Set retrieved user
	return authToken, nil
}

// RegenerateLogin get the refresh-token and returns a new one.
func (serv AuthService) RegenerateLogin(refreshToken string) (entities.AuthenticationToken, error) {
	authentication, err := serv.tokenRepository.CheckAuthenticationByRefreshToken(refreshToken)
	if err != nil || authentication.User.ID.IsEmpty() ||
		time.Now().After(authentication.ExpiresAt) {
		return entities.AuthenticationToken{}, autherrs.ErrInvalidAuthToken
	}
	// Update refresh token and expiration only (JWT will be generated in controller)
	if err = authentication.Regenerate(refreshTokenDuration); err != nil {
		return entities.AuthenticationToken{}, err
	}

	// Save AuthenticationToken update
	if err = serv.tokenRepository.UpsertAuthToken(&authentication); err != nil {
		return entities.AuthenticationToken{}, autherrs.ErrUpdateAuthToken
	}
	return authentication, nil
}

func (serv AuthService) VerifyUserCode(code string) error {
	user, err := serv.userRepository.GetUserByVerificationCode(code)
	if err != nil || user.ID.IsEmpty() {
		return autherrs.ErrVerificationCode
	}
	if err = serv.userRepository.SetUserVerified(user.ID); err != nil {
		return fmt.Errorf("failed to set user verified: %w", err)
	}
	return nil
}

func (serv AuthService) GeneratePasswordRecovery(userEmail string) error {
	user, err := serv.userRepository.GetUserByEmail(userEmail)
	if err != nil || user.ID.IsEmpty() {
		return autherrs.ErrUserEmailNotFound
	}

	var recoveryCode string
	if recoveryCode, err = authcore.GenerateRecoveryCode(userEmail); err != nil {
		return autherrs.ErrGenRecoveryCode
	}
	if err = serv.userRepository.SetRecoveryCode(user.ID, recoveryCode, time.Now().Add(resetCodeExpiration)); err != nil {
		return autherrs.ErrGenRecoveryCode
	}
	serv.mailerService.SendPasswordRecoveryEmail(user.Email, recoveryCode)
	return nil
}

func (serv AuthService) CheckRecoveryCode(
	identifier, recoveryCode string,
) (userID entities.HexID, err error) {
	var user entities.User
	user, err = serv.userRepository.GetUserByRecovery(
		identifier, recoveryCode, time.Now().Add(-resetCodeExpiration),
	)
	if err != nil || user.ID.IsEmpty() {
		return entities.HexID{}, autherrs.ErrUserRecoveryNotFound
	}

	return user.ID, nil
}

func (serv AuthService) ResetPassword(resetUser entities.UserBasic, recoveryCode string) error {
	userID, err := serv.CheckRecoveryCode(resetUser.Email, recoveryCode)
	if err != nil {
		return err
	}

	encryptedPassword, encryptErr := resetUser.EncryptPassword()
	if encryptErr != nil {
		return autherrs.ErrPasswordEncryption
	}
	if err = serv.userRepository.UpdatePassword(userID, string(encryptedPassword)); err != nil {
		return autherrs.ErrUpdatePassword
	}
	_ = serv.userRepository.SetRecoveryCode(userID, "", time.Time{})
	return nil
}

func (serv AuthService) GetObfuscatedEmail(username string) string {
	user, err := serv.userRepository.GetUserByUsername(username)
	if err != nil || user.ID.IsEmpty() {
		return ""
	}
	return user.ObfuscateEmail()
}
