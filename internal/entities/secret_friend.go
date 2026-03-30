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
	Timestamp

	ID              HexID
	Name            string
	Datetime        time.Time
	Location        string
	OwnerID         HexID
	InviteCode      string
	MaxDenyListSize uint8
	Status          SecretFriendStatus
	Participants    []Participant
}

func (sf *SecretFriend) Normalize() {
	if !sf.Datetime.IsZero() {
		sf.Datetime = sf.Datetime.In(time.UTC)
	}
}
