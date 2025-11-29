package entities

import "time"

type WishlistItem struct {
	ID            HexID
	ParticipantID HexID
	Label         string
	Comments      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
