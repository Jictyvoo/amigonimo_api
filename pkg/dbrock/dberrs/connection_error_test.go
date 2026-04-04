package dberrs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestErrDatabaseConnectionError(t *testing.T) {
	inner := errors.New("connection refused")

	tests := []struct {
		name  string
		inner error
		want  string
	}{
		{
			name: "no inner err",
			want: "database connection failed",
		},
		{
			name:  "with inner err",
			inner: inner,
			want:  "database connection failed: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrs.NewErrDatabaseConnection(tt.inner)
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrDatabaseConnectionStatusCode(t *testing.T) {
	err := dberrs.NewErrDatabaseConnection(nil)
	if got := err.StatusCode(); got != http.StatusServiceUnavailable {
		t.Errorf("StatusCode() = %d, want %d", got, http.StatusServiceUnavailable)
	}
}
