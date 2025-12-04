package drawresultctrl

// WishlistItem represents a wishlist item in the draw result
type WishlistItem struct {
	ItemID   string `json:"itemId"`
	Label    string `json:"label"`
	Comments string `json:"comments,omitempty"`
}

// DrawResultResponse represents the result of the draw for the current user
type DrawResultResponse struct {
	TargetUserID string         `json:"targetUserId"`
	TargetName   string         `json:"targetName"`
	Wishlist     []WishlistItem `json:"wishlist"`
}
