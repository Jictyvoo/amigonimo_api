package dberrs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestErrDatabaseValidationError(t *testing.T) {
	inner := errors.New("value out of range")

	tests := []struct {
		name  string
		field string
		inner error
		want  string
	}{
		{
			name: "no field no inner err",
			want: "database validation error",
		},
		{
			name:  "no field with inner err",
			inner: inner,
			want:  "database validation error: value out of range",
		},
		{
			name:  "field without inner err",
			field: "email",
			want:  "database validation error, field 'email'",
		},
		{
			name:  "field with inner err",
			field: "email",
			inner: inner,
			want:  "database validation error, field 'email': value out of range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrs.NewErrDatabaseValidation(tt.field, tt.inner)
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrDatabaseValidationStatusCode(t *testing.T) {
	err := dberrs.NewErrDatabaseValidation("email", nil)
	if got := err.StatusCode(); got != http.StatusBadRequest {
		t.Errorf("StatusCode() = %d, want %d", got, http.StatusBadRequest)
	}
}
