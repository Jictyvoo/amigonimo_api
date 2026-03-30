package entities

import (
	"strings"
	"time"
)

type Participant struct {
	Timestamp

	ID             HexID
	RelatedUser    User
	SecretFriendID HexID
	JoinedAt       time.Time
	IsReady        bool
}

func NewParticipant(secretFriendID HexID, relatedUser User) Participant {
	return Participant{SecretFriendID: secretFriendID, RelatedUser: relatedUser}
}

type WishlistItem struct {
	Timestamp

	ID       HexID
	Label    string
	Comments string
}

func (wi *WishlistItem) Normalize() {
	wi.Comments = strings.TrimSpace(wi.Comments)
	wi.Label = strings.TrimSpace(wi.Label)
	wi.Timestamp.Normalize()
}
