package entities

import (
	"strings"
	"time"
)

type DeniedUser struct {
	ID          HexID
	DeniedUsers Participant
	Timestamp
}

type Participant struct {
	ID             HexID
	RelatedUser    User
	SecretFriendID HexID
	JoinedAt       time.Time
	DenyList       []DeniedUser
	Wishlist       Wishlist
	Timestamp
}

type (
	WishlistItem struct {
		ID       HexID
		Label    string
		Comments string
		Timestamp
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
