package entities

import (
	"time"
)

type DrawResult struct {
	ID                    HexID
	SecretFriendID        HexID
	GiverParticipantID    HexID
	ReceiverParticipantID HexID
	CreatedAt             time.Time
	UpdatedAt             time.Time
	// Receiver info
	ReceiverUserID   HexID
	ReceiverName     string
	ReceiverEmail    string
	ReceiverWishlist []WishlistItem
}
