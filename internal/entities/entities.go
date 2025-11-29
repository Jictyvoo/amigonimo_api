package entities

import (
	"github.com/google/uuid"
)

type (
	Identifier uint64    // Integer based identifier
	HexID      uuid.UUID // Hex based id (like uuid)
)

func (id HexID) IsEmpty() bool {
	return id == HexID(uuid.Nil)
}

func (id HexID) String() string {
	return uuid.UUID(id).String()
}

func NewHexID() HexID {
	return HexID(uuid.New())
}
