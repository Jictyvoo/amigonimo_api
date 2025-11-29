package entities

import (
	"time"
)

type Participant struct {
	ID             HexID
	RelatedUser    User
	SecretFriendID HexID
	JoinedAt       time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
