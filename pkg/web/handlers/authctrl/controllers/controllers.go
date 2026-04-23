package controllers

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type AuthServiceFactory func(ctx context.Context) (authserv.AuthService, error)

type AuthenticationController struct {
	web.DefaultController

	secretKey       *rsa.PrivateKey
	authServFactory AuthServiceFactory
}

func NewAuthController(
	secretKey *rsa.PrivateKey,
	authServFactory AuthServiceFactory,
) AuthenticationController {
	return AuthenticationController{
		authServFactory: authServFactory,
		secretKey:       secretKey,
	}
}

func (h *AuthenticationController) authService(
	ctx context.Context,
) (authserv.AuthService, error) {
	return h.authServFactory(ctx)
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
