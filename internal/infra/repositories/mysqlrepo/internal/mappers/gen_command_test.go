package mappers_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func TestHexIDFromBytes(t *testing.T) {
	valid16 := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	over16 := make([]byte, 20)
	for i := range over16 {
		over16[i] = byte(i + 1)
	}

	tests := []struct {
		name  string
		input []byte
		want  entities.HexID
	}{
		{"valid 16 bytes", valid16, entities.HexID(valid16)},
		{"nil slice returns uuid.Nil", nil, entities.HexID(uuid.Nil)},
		{"fewer than 16 bytes returns uuid.Nil", []byte{1, 2, 3}, entities.HexID(uuid.Nil)},
		{"more than 16 bytes returns uuid.Nil", over16, entities.HexID(uuid.Nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mappers.HexIDFromBytes(tt.input))
		})
	}
}

func TestTimeFromNullTime(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name  string
		input sql.NullTime
		want  time.Time
	}{
		{"valid NullTime returns time", sql.NullTime{Time: now, Valid: true}, now},
		{"invalid NullTime returns zero", sql.NullTime{Time: now, Valid: false}, time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mappers.TimeFromNullTime(tt.input))
		})
	}
}

func TestTimeToNullTime(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name      string
		input     time.Time
		wantValid bool
		wantTime  time.Time
	}{
		{"non-zero time produces valid NullTime", now, true, now},
		{"zero time produces invalid NullTime", time.Time{}, false, time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mappers.TimeToNullTime(tt.input)
			assert.Equal(t, tt.wantValid, got.Valid)
			assert.Equal(t, tt.wantTime, got.Time)
		})
	}
}

func TestStringFromNullString(t *testing.T) {
	tests := []struct {
		name  string
		input sql.NullString
		want  string
	}{
		{"valid NullString returns string", sql.NullString{String: "hello", Valid: true}, "hello"},
		{"invalid NullString returns empty", sql.NullString{String: "ignored", Valid: false}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, mappers.StringFromNullString(tt.input))
		})
	}
}

func TestUUIDFromNullString(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name      string
		input     sql.NullString
		wantValid bool
		wantUUID  uuid.UUID
	}{
		{
			name:      "valid UUID string",
			input:     sql.NullString{String: validUUID.String(), Valid: true},
			wantValid: true,
			wantUUID:  validUUID,
		},
		{
			name:      "malformed UUID string returns invalid",
			input:     sql.NullString{String: "not-a-uuid", Valid: true},
			wantValid: false,
		},
		{
			name:      "invalid NullString returns invalid",
			input:     sql.NullString{Valid: false},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mappers.UUIDFromNullString(tt.input)
			assert.Equal(t, tt.wantValid, got.Valid)
			if tt.wantValid {
				assert.Equal(t, tt.wantUUID, got.UUID)
			}
		})
	}
}

func TestCopyTime(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name  string
		input time.Time
	}{
		{"non-zero time is preserved", now},
		{"zero time is preserved", time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.input, mappers.CopyTime(tt.input))
		})
	}
}
