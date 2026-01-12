package denylistctrl

// DeniedUserResponse represents a denied user in the deny list.
type DeniedUserResponse struct {
	UserID   string `json:"userId"`
	Fullname string `json:"fullname"`
}

// AddDenyListRequest represents the request to add a user to the deny list.
type AddDenyListRequest struct {
	TargetUserID string `json:"targetUserId" validate:"required"`
}
