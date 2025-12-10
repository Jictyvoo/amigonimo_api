package authserv

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type MailerService interface {
	SendActivationEmail(email string, verificationToken string)
	SendPasswordRecoveryEmail(email string, recoveryCode string)
}

//go:generate mockgen -destination=../../mocks/token_repo_mock.go -package=mocks github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv TokenRepository

type (
	UserRetriever interface {
		GetUserByAuthToken(token string) (entities.User, error)
	}
	TokenRepository interface {
		UserRetriever
		GetAuthenticationToken(userID entities.HexID) (entities.AuthenticationToken, error)
		UpdateUserAuthToken(authentication *entities.AuthenticationToken) error
		CheckAuthenticationByRefreshToken(authToken string) (*entities.AuthenticationToken, error)
	}
)

//go:generate mockgen -destination=../../mocks/userauth_repo_mock.go -package=mocks github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv UserAuthRepository

type UserAuthRepository interface {
	GetUserByUsername(username string) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetUserByVerificationCode(code string) (entities.User, error)
	GetUserByRecovery(userEmail string, code string, expiredAt time.Time) (entities.User, error)
	CreateUser(user entities.User, token string) error
	SetUserVerified(userID entities.HexID) error
	SetRecoveryCode(userID entities.HexID, code string) error
	UpdatePassword(userID entities.HexID, newPassword string) error
	GetUserByEmailOrUsername(email, username string) (entities.User, error)
}
