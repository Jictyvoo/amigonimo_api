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

type (
	DrawResultItem struct {
		Timestamp

		Giver    Participant
		Receiver Participant
	}
	DrawResult struct {
		GiverReceivers []DrawResultItem
	}
)

type SecretFriend struct {
	Timestamp

	ID              HexID
	Name            string
	Datetime        *time.Time
	Location        string
	OwnerID         HexID
	InviteCode      string
	MaxDenyListSize uint8
	Status          SecretFriendStatus
	Participants    []Participant
	DrawResult      *DrawResult
}

func (sf *SecretFriend) Normalize() {
	if sf.Datetime != nil {
		*sf.Datetime = sf.Datetime.In(time.UTC)
	}
}
