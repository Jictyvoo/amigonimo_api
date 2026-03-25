package entities

import (
	"testing"
	"time"
)

func TestBasicAuthTokenRegenerate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "generates token fields for positive duration",
			duration: 30 * time.Minute,
		},
		{
			name:     "still generates token fields for zero duration",
			duration: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var token BasicAuthToken
			before := time.Now()

			err := token.Regenerate(tt.duration)
			if err != nil {
				t.Fatalf("Regenerate() error = %v, want nil", err)
			}
			if token.AuthToken == "" {
				t.Fatal("Regenerate() left AuthToken empty")
			}
			if !token.RefreshToken.Valid {
				t.Fatal("Regenerate() left RefreshToken invalid")
			}
			if token.RefreshToken.UUID.String() == "" {
				t.Fatal("Regenerate() left RefreshToken UUID empty")
			}

			minExpiration := before.Add(tt.duration)
			if token.ExpiresAt.Before(minExpiration) {
				t.Fatalf("ExpiresAt = %v, want >= %v", token.ExpiresAt, minExpiration)
			}
		})
	}
}
