package jwtware_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jictyvoo/amigonimo_api/pkg/web/middlewares/jwtware"
)

func TestNormalizeConfig_PanicsWithoutKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when neither SigningKey nor KeyFunc is set")
		}
	}()
	// Empty config — should panic.
	_ = jwtware.New[jwt.MapClaims, *jwt.MapClaims]()
}

func TestNormalizeConfig_DefaultSchemeIsBearer(t *testing.T) {
	handler := jwtware.New[jwt.MapClaims, *jwt.MapClaims](
		jwtware.MapClaimsConfig{
			SigningKey: jwtware.SigningKey{Key: []byte(testHMACSecret)},
		},
	)

	// Build a request with a valid Bearer token just to verify the middleware
	// passes through — the default scheme must be Bearer.
	token := makeHMACToken(t, validClaims())
	req, rec := newRequest(t, "Bearer "+token)
	handler(okHandler(new(false))).ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d (default scheme should be Bearer)", rec.Code)
	}
}

func TestNormalizeConfig_CustomScheme(t *testing.T) {
	handler := jwtware.New[jwt.MapClaims, *jwt.MapClaims](
		jwtware.MapClaimsConfig{
			SigningKey: jwtware.SigningKey{Key: []byte(testHMACSecret)},
			Scheme:     jwtware.SchemeBasic,
		},
	)

	token := makeHMACToken(t, validClaims())
	req, rec := newRequest(t, "Basic "+token)
	handler(okHandler(new(false))).ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200 with Basic scheme, got %d", rec.Code)
	}
}

func TestNormalizeConfig_CustomKeyFunc(t *testing.T) {
	customKeyFuncCalled := false
	keyFunc := func(token *jwt.Token) (any, error) {
		customKeyFuncCalled = true
		return []byte(testHMACSecret), nil
	}

	handler := jwtware.New[jwt.MapClaims, *jwt.MapClaims](
		jwtware.MapClaimsConfig{
			KeyFunc: keyFunc,
		},
	)

	token := makeHMACToken(t, validClaims())
	req, rec := newRequest(t, "Bearer "+token)
	handler(okHandler(new(false))).ServeHTTP(rec, req)

	if !customKeyFuncCalled {
		t.Error("custom KeyFunc was not called")
	}
	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
