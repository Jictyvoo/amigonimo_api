package entities

import (
	"strings"
	"time"
)

type DeniedUser struct {
	Timestamp

	ID               HexID
	InnerParticipant Participant
}

type Participant struct {
	Timestamp

	ID             HexID
	RelatedUser    User
	SecretFriendID HexID
	JoinedAt       time.Time
	DenyList       []DeniedUser
	Wishlist       Wishlist
	IsReady        bool
}

func NewParticipant(secretFriendID HexID, relatedUser User) Participant {
	return Participant{SecretFriendID: secretFriendID, RelatedUser: relatedUser}
}

type (
	WishlistItem struct {
		Timestamp

		ID       HexID
		Label    string
		Comments string
	}

	Wishlist struct {
		Items []WishlistItem
	}
)

func (wi *WishlistItem) Normalize() {
	wi.Comments = strings.TrimSpace(wi.Comments)
	wi.Label = strings.TrimSpace(wi.Label)
	wi.Timestamp.Normalize()
}
