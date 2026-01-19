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

func NewHexID() (HexID, error) {
	uid, err := uuid.NewV7()
	return HexID(uid), err
}

func ParseHexID(input string) (HexID, error) {
	loadId, err := uuid.Parse(input)
	if err != nil {
		return HexID{}, err
	}
	return HexID(loadId), nil
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
