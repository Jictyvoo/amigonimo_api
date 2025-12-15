package jwtware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func New[T any, V TokenClaims[T]](optConfig ...Config[T, V]) func(http.Handler) http.Handler {
	conf := normalizeConfig(optConfig...)
	return func(next http.Handler) http.Handler {
		ware := JWTWare[T, V]{
			conf: &conf,
			next: next,
		}
		return http.HandlerFunc(ware.handlerFunc)
	}
}

type (
	ctxKeyAuthClaims                 struct{}
	JWTWare[T any, V TokenClaims[T]] struct {
		conf *Config[T, V]
		next http.Handler
	}
)

func (jw JWTWare[T, V]) handlerFunc(w http.ResponseWriter, r *http.Request) {
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

func ClaimsFromContext[T jwt.Claims](ctx context.Context) (tkClaims T, err error) {
	var ok bool
	if tkClaims, ok = ctx.Value(ctxKeyAuthClaims{}).(T); !ok {
		return tkClaims, errors.New("user ID not found in context")
	}
	return tkClaims, nil
}
