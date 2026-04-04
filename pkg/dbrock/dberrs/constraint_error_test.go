package dberrs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestErrDatabaseConstraintError(t *testing.T) {
	inner := errors.New("duplicate entry")

	tests := []struct {
		name       string
		constraint string
		inner      error
		want       string
	}{
		{
			name: "no constraint no inner err",
			want: "database constraint violation",
		},
		{
			name:  "no constraint with inner err",
			inner: inner,
			want:  "database constraint violation: duplicate entry",
		},
		{
			name:       "constraint with inner err",
			constraint: "users_email_unique",
			inner:      inner,
			want:       "database constraint violation, constraint 'users_email_unique': duplicate entry",
		},
		{
			// constraint field is only shown when there is a wrapped error
			name:       "constraint without inner err",
			constraint: "users_email_unique",
			want:       "database constraint violation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrs.NewErrDatabaseConstraint(tt.constraint, tt.inner)
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrDatabaseConstraintStatusCode(t *testing.T) {
	err := dberrs.NewErrDatabaseConstraint("unique", nil)
	if got := err.StatusCode(); got != http.StatusConflict {
		t.Errorf("StatusCode() = %d, want %d", got, http.StatusConflict)
	}
}
