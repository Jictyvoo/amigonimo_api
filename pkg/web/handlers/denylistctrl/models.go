package denylistctrl

// AddDenyListRequest represents the request to add a user to the deny list.
type AddDenyListRequest struct {
	TargetUserID string `json:"targetUserId" validate:"required"`
}

// DeniedUserResponse represents a denied user in the deny list.
type DeniedUserResponse struct {
	UserID   string `json:"userId"`
	Fullname string `json:"fullname"`
}

// RemoveDenyListEntryResponse represents the response after deleting a deny list entry.
type RemoveDenyListEntryResponse struct {
	Success   bool   `json:"success"`
	DeletedID string `json:"deletedId"`
}
