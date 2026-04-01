package drawdto

import "github.com/jictyvoo/amigonimo_api/internal/entities"

// WishlistEntry is a minimal read-model for a single wishlist item
// as surfaced by the draw result. It carries only the display fields.
type WishlistEntry struct {
	Label    string
	Comments string
}

// DrawResultItem is the draw pipeline DTO.
// The execute path only fills the participant ID fields.
// The getresult path additionally fills the display fields.
type DrawResultItem struct {
	GiverParticipantID    entities.HexID
	ReceiverParticipantID entities.HexID
	// Display fields — populated by the getresult read path only.
	ReceiverUserID   entities.HexID
	ReceiverEmail    string
	ReceiverFullName string
	ReceiverWishlist []WishlistEntry
}
