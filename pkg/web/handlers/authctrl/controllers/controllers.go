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
	secretKey []byte, authServ authserv.AuthService,
) AuthenticationController {
	return AuthenticationController{
		authServ:  authServ,
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
