package entities

import "time"

type UserProfile struct {
	Timestamp

	ID        HexID
	UserID    HexID
	FullName  string
	Nickname  string
	ImageLink string
	Birthday  time.Time
	Address   string
}
