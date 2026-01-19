package jwtware

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims[T any] interface {
	*T
	jwt.Claims
}

type AuthScheme string

const (
	SchemeBearer AuthScheme = "Bearer"
	SchemeBasic  AuthScheme = "Basic"
)

// SigningKey holds information about the recognized cryptographic keys used to sign JWTs by this program.
type SigningKey struct {
	// JWTAlg is the algorithm used to sign JWTs. If this value is a non-empty string, this will be checked against the
	// "alg" value in the JWT header.
	//
	// https://www.rfc-editor.org/rfc/rfc7518#section-3.1
	JWTAlg string
	// Key is the cryptographic key used to sign JWTs. For supported types, please see
	// https://github.com/golang-jwt/jwt.
	Key any
}

type (
	MapClaimsConfig                 = Config[jwt.MapClaims, *jwt.MapClaims]
	Config[T any, V TokenClaims[T]] struct {
		// SigningKey is the primary key used to validate tokens.
		SigningKey SigningKey

		// KeyFunc provides the public key for JWT verification.
		// It handles algorithm verification and key selection.
		// At least one of the following is required: KeyFunc or SigningKey.
		KeyFunc jwt.Keyfunc

		// TokenProcessorFunc processes the token extracted.
		// Optional. Default: nil
		TokenProcessorFunc func(token string, keyFunc jwt.Keyfunc) (T, error)

		// ErrorHandler deals with all errors raised.
		// Optional. Default: nil
		ErrorHandler func(w http.ResponseWriter, err error)

		// Scheme represents the desired scheme to check on Authorization header.
		// Optional. Default: SchemeBearer
		Scheme AuthScheme
	}
)

func normalizeConfig[T any, V TokenClaims[T]](optConfigList ...Config[T, V]) (conf Config[T, V]) {
	if len(optConfigList) > 0 {
		conf = optConfigList[0]
	}

	if conf.SigningKey.Key == nil && conf.KeyFunc == nil {
		panic("JWT secret key is required")
	}

	if conf.TokenProcessorFunc == nil {
		conf.TokenProcessorFunc = defaultTokenProcessor[T, V]
	}

	if conf.KeyFunc == nil {
		conf.KeyFunc = signingKeyFunc(conf.SigningKey)
	}

	if conf.Scheme == "" {
		conf.Scheme = SchemeBearer
	}
	return conf
}

func defaultTokenProcessor[T any, V TokenClaims[T]](token string, keyFunc jwt.Keyfunc) (T, error) {
	var storedClaims T
	decodedToken, err := jwt.ParseWithClaims(token, V(&storedClaims), keyFunc)
	if err != nil {
		return storedClaims, err
	}

	if !decodedToken.Valid {
		return storedClaims, errors.New("invalid token")
	}

	return storedClaims, nil
}
