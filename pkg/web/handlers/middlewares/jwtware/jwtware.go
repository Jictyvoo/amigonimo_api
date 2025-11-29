package jwtware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func New[T jwt.Claims](optConfig ...Config[T]) func(http.Handler) http.Handler {
	conf := normalizeConfig(optConfig...)
	return func(next http.Handler) http.Handler {
		ware := JWTWare[T]{
			conf: &conf,
			next: next,
		}
		return http.HandlerFunc(ware.handlerFunc)
	}
}

type (
	ctxKeyAuthClaims      struct{}
	JWTWare[T jwt.Claims] struct {
		conf *Config[T]
		next http.Handler
	}
)

func (jw JWTWare[T]) handlerFunc(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>" format
	authScheme, token, _ := strings.Cut(authHeader, " ")
	if authScheme != string(jw.conf.Scheme) {
		http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
		return
	}

	tkClaims, err := jw.conf.TokenProcessorFunc(token, jw.conf.KeyFunc)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Store claims in context
	ctx := context.WithValue(r.Context(), ctxKeyAuthClaims{}, tkClaims)
	jw.next.ServeHTTP(w, r.WithContext(ctx))
}

func GetClaimsFromContext(ctx context.Context) (jwt.Claims, error) {
	tkClaims, ok := ctx.Value(ctxKeyAuthClaims{}).(jwt.Claims)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}
	return tkClaims, nil
}
