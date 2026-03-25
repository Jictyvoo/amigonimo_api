package entities

import (
	"testing"
	"time"
)

func TestSecretFriendNormalize(t *testing.T) {
	t.Parallel()

	localTZ := time.FixedZone("UTC-3", -3*60*60)
	localTime := time.Date(2026, time.March, 25, 20, 15, 0, 0, localTZ)

	tests := []struct {
		name      string
		input     SecretFriend
		assertion func(*testing.T, SecretFriend)
	}{
		{
			name: "converts datetime to utc",
			input: SecretFriend{
				Datetime: localTime,
			},
			assertion: func(t *testing.T, got SecretFriend) {
				t.Helper()
				if got.Datetime.Location() != time.UTC {
					t.Fatalf("Datetime location = %v, want UTC", got.Datetime.Location())
				}
				if !got.Datetime.Equal(localTime.UTC()) {
					t.Fatalf("Datetime = %v, want %v", got.Datetime, localTime.UTC())
				}
			},
		},
		{
			name:  "keeps zero datetime unchanged",
			input: SecretFriend{},
			assertion: func(t *testing.T, got SecretFriend) {
				t.Helper()
				if !got.Datetime.IsZero() {
					t.Fatalf("Datetime = %v, want zero", got.Datetime)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := tt.input
			sf.Normalize()
			tt.assertion(t, sf)
		})
	}
}
