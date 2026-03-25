package web

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
)

func TestMapError(t *testing.T) {
	internalErr := errors.New("boom")
	appErr := apperr.Invalid("invalid_input", "invalid input", internalErr)
	authErr := autherrs.NewErrLogin(internalErr)
	httpErr := &fuego.HTTPError{
		Err:    errors.New("bad request"),
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: "bad request",
	}

	tests := []struct {
		name         string
		err          error
		wantStatus   int
		wantDetail   string
		wantType     string
		wantSameHTTP bool
	}{
		{
			name:       "app error",
			err:        appErr,
			wantStatus: http.StatusBadRequest,
			wantDetail: "invalid input",
			wantType:   "invalid_input",
		},
		{
			name:       "auth error",
			err:        authErr,
			wantStatus: http.StatusInternalServerError,
			wantDetail: "failed to log in",
			wantType:   "auth_login_failed",
		},
		{
			name:         "http error passthrough",
			err:          httpErr,
			wantStatus:   http.StatusBadRequest,
			wantDetail:   "bad request",
			wantSameHTTP: true,
		},
		{
			name:       "unknown error becomes internal server error",
			err:        errors.New("unknown"),
			wantStatus: http.StatusInternalServerError,
			wantDetail: http.StatusText(http.StatusInternalServerError),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				gotErr := MapError(tt.err)
				if gotErr == nil {
					t.Fatal("MapError() = nil, want error")
				}

				gotHTTP, ok := errors.AsType[*fuego.HTTPError](gotErr)
				if !ok {
					t.Fatalf("MapError() type = %T, want *fuego.HTTPError", gotErr)
				}

				if tt.wantSameHTTP && !errors.Is(gotHTTP, httpErr) {
					t.Fatal("MapError() should have returned the original HTTP error")
				}
				if gotHTTP.Status != tt.wantStatus {
					t.Fatalf("Status = %d, want %d", gotHTTP.Status, tt.wantStatus)
				}
				if gotHTTP.Detail != tt.wantDetail {
					t.Fatalf("Detail = %q, want %q", gotHTTP.Detail, tt.wantDetail)
				}
				if gotHTTP.Type != tt.wantType {
					t.Fatalf("Type = %q, want %q", gotHTTP.Type, tt.wantType)
				}
			},
		)
	}
}
