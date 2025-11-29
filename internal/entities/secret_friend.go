package entities

import (
	"time"
)

type SecretFriendStatus string

const (
	StatusDraft  SecretFriendStatus = "draft"
	StatusOpen   SecretFriendStatus = "open"
	StatusDrawn  SecretFriendStatus = "drawn"
	StatusClosed SecretFriendStatus = "closed"
)

type SecretFriend struct {
	ID                HexID
	Name              string
	Datetime          *time.Time
	Location          string
	OwnerID           HexID
	InviteCode        string
	InviteLink        string
	Status            SecretFriendStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ParticipantsCount int
}

func (sf *SecretFriend) Normalize() {
	if sf.Datetime != nil {
		*sf.Datetime = sf.Datetime.In(time.UTC)
	}
}
