package dberrs_test

import (
	"errors"
	"testing"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestErrDatabaseTransactionError(t *testing.T) {
	inner := errors.New("tx failed")

	tests := []struct {
		name      string
		operation string
		inner     error
		want      string
	}{
		{
			name: "no operation no inner err",
			want: "database transaction error",
		},
		{
			name:  "no operation with inner err",
			inner: inner,
			want:  "database transaction error: tx failed",
		},
		{
			name:      "operation without inner err",
			operation: "commit",
			want:      "database transaction erroroperation 'commit'",
		},
		{
			name:      "operation with inner err",
			operation: "commit",
			inner:     inner,
			want:      "database transaction erroroperation 'commit': tx failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dberrs.NewErrDatabaseTransaction(tt.operation, tt.inner)
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}
