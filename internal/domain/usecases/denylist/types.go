package denylist

import "github.com/jictyvoo/amigonimo_api/internal/entities"

// DeniedEntry is a read-model DTO returned by denylist queries.
// It carries the minimal display data needed by the handler.
type DeniedEntry struct {
	entities.Timestamp
	ID           entities.HexID
	DeniedUserID entities.HexID
	Username     string
	Email        string
	FullName     string
}
