package controllers

import (
	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
)

type AuthenticationController struct {
	secretKey []byte
	authServ  authserv.AuthService
}

func NewAuthController(
	secretKey []byte, mailerService authserv.MailerService,
	userRepository authserv.UserAuthRepository, tokenRepository authserv.TokenRepository,
) AuthenticationController {
	return AuthenticationController{
		authServ:  authserv.NewAuthService(userRepository, tokenRepository, mailerService),
		secretKey: secretKey,
	}
}

// NewHTTPError creates a new HTTP error with the given status code and message
func NewHTTPError(statusCode int, message string) *fuego.HTTPError {
	return &fuego.HTTPError{
		Status: statusCode,
		Detail: message,
	}
}
