package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	Identifier uint64    // Integer based identifier
	HexID      uuid.UUID // Hex based id (like uuid)

	Timestamp struct {
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func NewHexID() HexID {
	return HexID(uuid.New())
}

func (id HexID) IsEmpty() bool {
	return id == HexID(uuid.Nil)
}

func (id HexID) String() string {
	return uuid.UUID(id).String()
}

func (ts *Timestamp) Normalize() {
	if ts.CreatedAt.IsZero() {
		ts.CreatedAt = time.Now()
	}
	if ts.UpdatedAt.IsZero() {
		ts.UpdatedAt = ts.CreatedAt
	}
}
