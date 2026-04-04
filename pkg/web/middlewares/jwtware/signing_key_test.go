package jwtware_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jictyvoo/amigonimo_api/pkg/web/middlewares/jwtware"
)

func TestSigningKeyFunc_AlgorithmEnforcement(t *testing.T) {
	validToken := makeHMACToken(t, validClaims())

	tests := []struct {
		name       string
		signingKey jwtware.SigningKey
		wantStatus int
	}{
		{
			name: "correct algorithm passes",
			signingKey: jwtware.SigningKey{
				JWTAlg: jwt.SigningMethodHS256.Alg(),
				Key:    []byte(testHMACSecret),
			},
			wantStatus: 200,
		},
		{
			name: "wrong algorithm rejected",
			signingKey: jwtware.SigningKey{
				JWTAlg: "RS256", // token is HS256
				Key:    []byte(testHMACSecret),
			},
			wantStatus: 401,
		},
		{
			name: "no algorithm constraint accepts any",
			signingKey: jwtware.SigningKey{
				Key: []byte(testHMACSecret),
			},
			wantStatus: 200,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handler := jwtware.New[jwt.MapClaims, *jwt.MapClaims](
					jwtware.MapClaimsConfig{
						SigningKey: tt.signingKey,
					},
				)
				req, rec := newRequest(t, "Bearer "+validToken)
				handler(okHandler(new(false))).ServeHTTP(rec, req)
				if rec.Code != tt.wantStatus {
					t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
				}
			},
		)
	}
}

const testHMACSecret = "super-secret-key"

// makeHMACToken signs a MapClaims token with HS256 using testHMACSecret.
func makeHMACToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte(testHMACSecret))
	if err != nil {
		t.Fatalf("makeHMACToken: %v", err)
	}
	return signed
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
}

func expiredClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
}
