package dberrs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestErrDatabaseNotFoundError(t *testing.T) {
	inner := errors.New("sql: no rows")

	tests := []struct {
		name       string
		resource   string
		identifier string
		inner      error
		want       string
	}{
		{
			name: "no fields no inner err",
			want: "resource not found",
		},
		{
			name:  "no fields with inner err",
			inner: inner,
			want:  "resource not found: sql: no rows",
		},
		{
			name:       "resource and identifier no inner err",
			resource:   "user",
			identifier: "42",
			want:       "user with identifier '42' not found",
		},
		{
			name:       "resource and identifier with inner err",
			resource:   "user",
			identifier: "42",
			inner:      inner,
			want:       "user with identifier '42' not found: sql: no rows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrs.NewErrDatabaseNotFound(tt.resource, tt.identifier, tt.inner)
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrDatabaseNotFoundStatusCode(t *testing.T) {
	err := dberrs.NewErrDatabaseNotFound("user", "1", nil)
	if got := err.StatusCode(); got != http.StatusNotFound {
		t.Errorf("StatusCode() = %d, want %d", got, http.StatusNotFound)
	}
}
