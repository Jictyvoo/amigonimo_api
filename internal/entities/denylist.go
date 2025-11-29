package entities

import (
	"time"
)

type Denylist struct {
	ID            HexID
	ParticipantID HexID
	DeniedUsers   []Participant
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
