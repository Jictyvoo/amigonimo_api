package controllers

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type AuthenticationController struct {
	web.DefaultController

	secretKey *rsa.PrivateKey
	authServ  authserv.AuthService
}

func NewAuthController(
	secretKey *rsa.PrivateKey, authServ authserv.AuthService,
) AuthenticationController {
	return AuthenticationController{
		authServ:  authServ,
		secretKey: secretKey,
	}
}

// generateJWT creates a JWT token from AuthenticationToken with user info.
func (h *AuthenticationController) generateJWT(
	authToken entities.AuthenticationToken,
) (string, error) {
	var verifiedAt *int64
	if !authToken.User.VerifiedAt.IsZero() {
		verifiedAt = new(int64)
		*verifiedAt = authToken.User.VerifiedAt.Unix()
	}
	claims := jwt.MapClaims{
		"userID":     authToken.User.ID.String(),
		"username":   authToken.User.Username,
		"tokenId":    authToken.AuthToken,
		"verifiedAt": verifiedAt,
		"exp":        authToken.ExpiresAt.Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	return token.SignedString(h.secretKey)
}
