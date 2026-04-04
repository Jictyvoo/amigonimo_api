package jwtware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jictyvoo/amigonimo_api/pkg/web/middlewares/jwtware"
)

// newRequest builds a GET request and recorder with the given Authorization value.
func newRequest(t *testing.T, authHeader string) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	return req, httptest.NewRecorder()
}

// okHandler returns a handler that sets *called=true and writes 200.
func okHandler(called *bool) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			*called = true
			w.WriteHeader(http.StatusOK)
		},
	)
}

func TestJWTWareHandlerFunc(t *testing.T) {
	validToken := makeHMACToken(t, validClaims())
	expiredToken := makeHMACToken(t, expiredClaims())

	tests := []struct {
		name       string
		authHeader string
		signingKey []byte
		wantStatus int
		wantNext   bool
	}{
		{
			name:       "missing Authorization header returns 401",
			authHeader: "",
			signingKey: []byte(testHMACSecret),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong scheme returns 401",
			authHeader: "Basic " + validToken,
			signingKey: []byte(testHMACSecret),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed token returns 401",
			authHeader: "Bearer not.a.valid.token",
			signingKey: []byte(testHMACSecret),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "expired token returns 401",
			authHeader: "Bearer " + expiredToken,
			signingKey: []byte(testHMACSecret),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong signing key returns 401",
			authHeader: "Bearer " + validToken,
			signingKey: []byte("completely-wrong-key"),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid token calls next handler",
			authHeader: "Bearer " + validToken,
			signingKey: []byte(testHMACSecret),
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				nextCalled := false
				handler := jwtware.New[jwt.MapClaims, *jwt.MapClaims](
					jwtware.MapClaimsConfig{
						SigningKey: jwtware.SigningKey{Key: tt.signingKey},
					},
				)(okHandler(&nextCalled))

				req, rec := newRequest(t, tt.authHeader)
				handler.ServeHTTP(rec, req)

				if rec.Code != tt.wantStatus {
					t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
				}
				if nextCalled != tt.wantNext {
					t.Errorf("nextCalled = %v, want %v", nextCalled, tt.wantNext)
				}
			},
		)
	}
}

func TestClaimsFromContext(t *testing.T) {
	t.Run(
		"present claims are returned", func(t *testing.T) {
			var capturedCtx context.Context
			handler := jwtware.New[jwt.MapClaims, *jwt.MapClaims](
				jwtware.MapClaimsConfig{
					SigningKey: jwtware.SigningKey{Key: []byte(testHMACSecret)},
				},
			)(
				http.HandlerFunc(
					func(_ http.ResponseWriter, r *http.Request) {
						capturedCtx = r.Context()
					},
				),
			)

			token := makeHMACToken(t, validClaims())
			req, rec := newRequest(t, "Bearer "+token)
			handler.ServeHTTP(rec, req)

			claims, err := jwtware.ClaimsFromContext[jwt.MapClaims](capturedCtx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if claims["sub"] != "user-123" {
				t.Errorf("claims[sub] = %v, want %q", claims["sub"], "user-123")
			}
		},
	)

	t.Run(
		"absent claims return error", func(t *testing.T) {
			_, err := jwtware.ClaimsFromContext[jwt.MapClaims](context.Background())
			if err == nil {
				t.Error("expected error for empty context, got nil")
			}
		},
	)
}
