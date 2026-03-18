package apperr

import (
	"errors"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestFrom(t *testing.T) {
	infraErr := dberrs.NewErrDatabaseNotFound("secret_friend", "123", errors.New("missing"))
	baseErr := Conflict("existing_code", "existing message", errors.New("wrapped"))

	tests := []struct {
		name             string
		code             string
		publicMessage    string
		err              error
		wantNil          bool
		wantCode         string
		wantMessage      string
		wantStatusCode   int
		wantSameInstance bool
	}{
		{
			name:    "nil error",
			code:    "ignored",
			err:     nil,
			wantNil: true,
		},
		{
			name:             "reuse concrete app error when not overriding",
			err:              baseErr,
			wantCode:         "existing_code",
			wantMessage:      "existing message",
			wantStatusCode:   409,
			wantSameInstance: true,
		},
		{
			name:           "wrap app error with new code and message",
			code:           "new_code",
			publicMessage:  "new message",
			err:            baseErr,
			wantCode:       "new_code",
			wantMessage:    "new message",
			wantStatusCode: 409,
		},
		{
			name:           "translate infra not found",
			code:           "secret_friend_not_found",
			publicMessage:  "secret friend not found",
			err:            infraErr,
			wantCode:       "secret_friend_not_found",
			wantMessage:    "secret friend not found",
			wantStatusCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := From(tt.code, tt.publicMessage, tt.err)

				if tt.wantNil {
					if got != nil {
						t.Fatalf("From() = %#v, want nil", got)
					}
					return
				}

				if got == nil {
					t.Fatal("From() = nil, want error")
				}
				if got.Code() != tt.wantCode {
					t.Fatalf("Code() = %q, want %q", got.Code(), tt.wantCode)
				}
				if got.Error() != tt.wantMessage {
					t.Fatalf("Error() = %q, want %q", got.Error(), tt.wantMessage)
				}
				if got.StatusCode() != tt.wantStatusCode {
					t.Fatalf("StatusCode() = %d, want %d", got.StatusCode(), tt.wantStatusCode)
				}
				if tt.wantSameInstance && !errors.Is(got, baseErr) {
					t.Fatalf("From() should have reused the original error instance")
				}
			},
		)
	}
}
