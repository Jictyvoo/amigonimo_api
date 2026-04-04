package dberrs_test

import (
	"errors"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestErrDatabaseQueryError(t *testing.T) {
	inner := errors.New("syntax error near SELECT")

	tests := []struct {
		name  string
		query string
		inner error
		want  string
	}{
		{
			name: "no query no inner err",
			want: "database query error",
		},
		{
			name:  "no query with inner err",
			inner: inner,
			want:  "database query error: syntax error near SELECT",
		},
		{
			name:  "query without inner err",
			query: "SELECT * FROM users",
			want:  "database query error, query 'SELECT * FROM users'",
		},
		{
			name:  "query with inner err",
			query: "SELECT * FROM users",
			inner: inner,
			want:  "database query error, query 'SELECT * FROM users': syntax error near SELECT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrs.NewErrDatabaseQuery(tt.query, tt.inner)
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
