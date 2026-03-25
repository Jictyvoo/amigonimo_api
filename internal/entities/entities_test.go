package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewHexID(t *testing.T) {
	t.Parallel()

	id, err := NewHexID()
	if err != nil {
		t.Fatalf("NewHexID() error = %v, want nil", err)
	}
	if id.IsEmpty() {
		t.Fatal("NewHexID() returned empty id")
	}
}

func TestNewHexIDFromBytes(t *testing.T) {
	t.Parallel()

	validUUID := uuid.MustParse("0195d1f8-1f84-7a8a-a8fd-764e5e67ad11")

	tests := []struct {
		name    string
		input   []byte
		want    HexID
		wantErr bool
	}{
		{
			name:  "creates id from valid bytes",
			input: validUUID[:],
			want:  HexID(validUUID),
		},
		{
			name:    "returns error for invalid byte length",
			input:   []byte{1, 2, 3},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewHexIDFromBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewHexIDFromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("NewHexIDFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseHexID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    HexID
		wantErr bool
	}{
		{
			name:  "parses valid uuid string",
			input: "0195d1f8-1f84-7a8a-a8fd-764e5e67ad11",
			want:  HexID(uuid.MustParse("0195d1f8-1f84-7a8a-a8fd-764e5e67ad11")),
		},
		{
			name:    "returns error for invalid uuid string",
			input:   "not-a-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseHexID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseHexID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("ParseHexID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexIDHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		id        HexID
		wantEmpty bool
		wantStr   string
	}{
		{
			name:      "nil uuid is empty",
			id:        HexID(uuid.Nil),
			wantEmpty: true,
			wantStr:   uuid.Nil.String(),
		},
		{
			name:    "regular uuid is not empty",
			id:      HexID(uuid.MustParse("0195d1f8-1f84-7a8a-a8fd-764e5e67ad11")),
			wantStr: "0195d1f8-1f84-7a8a-a8fd-764e5e67ad11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.id.IsEmpty(); got != tt.wantEmpty {
				t.Fatalf("IsEmpty() = %v, want %v", got, tt.wantEmpty)
			}
			if got := tt.id.String(); got != tt.wantStr {
				t.Fatalf("String() = %q, want %q", got, tt.wantStr)
			}
		})
	}
}

func TestTimestampNormalize(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 25, 10, 30, 0, 0, time.UTC)
	later := now.Add(2 * time.Hour)

	tests := []struct {
		name          string
		timestamp     Timestamp
		assertionFunc func(*testing.T, Timestamp)
	}{
		{
			name:      "fills both timestamps when zero",
			timestamp: Timestamp{},
			assertionFunc: func(t *testing.T, got Timestamp) {
				t.Helper()
				if got.CreatedAt.IsZero() {
					t.Fatal("CreatedAt is zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Fatal("UpdatedAt is zero")
				}
				if !got.UpdatedAt.Equal(got.CreatedAt) {
					t.Fatalf(
						"UpdatedAt = %v, want equal to CreatedAt %v",
						got.UpdatedAt,
						got.CreatedAt,
					)
				}
			},
		},
		{
			name:      "keeps created at and fills updated at",
			timestamp: Timestamp{CreatedAt: now},
			assertionFunc: func(t *testing.T, got Timestamp) {
				t.Helper()
				if !got.CreatedAt.Equal(now) {
					t.Fatalf("CreatedAt = %v, want %v", got.CreatedAt, now)
				}
				if !got.UpdatedAt.Equal(now) {
					t.Fatalf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
				}
			},
		},
		{
			name:      "keeps both timestamps when already set",
			timestamp: Timestamp{CreatedAt: now, UpdatedAt: later},
			assertionFunc: func(t *testing.T, got Timestamp) {
				t.Helper()
				if !got.CreatedAt.Equal(now) {
					t.Fatalf("CreatedAt = %v, want %v", got.CreatedAt, now)
				}
				if !got.UpdatedAt.Equal(later) {
					t.Fatalf("UpdatedAt = %v, want %v", got.UpdatedAt, later)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := tt.timestamp
			ts.Normalize()
			tt.assertionFunc(t, ts)
		})
	}
}
