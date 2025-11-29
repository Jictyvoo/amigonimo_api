package jwtware

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidSigningKey    = errors.New("invalid signing key")
	ErrInvalidToken         = errors.New("provided token is invalid")
	ErrInvalidAuthorization = errors.New("invalid authorization header format")
	ErrMissingHeader        = errors.New("missing authorization header")
)

func signingKeyFunc(key SigningKey) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		if key.JWTAlg != "" {
			algorithm, _ := token.Header["alg"].(string)
			switch algorithm {
			case key.JWTAlg:
				// Correct, do nothing
			default:
				return nil, fmt.Errorf(
					"unexpected jwt signing method: expected: %q: got: %q",
					key.JWTAlg, algorithm,
				)
			}
		}
		return key.Key, nil
	}
}
