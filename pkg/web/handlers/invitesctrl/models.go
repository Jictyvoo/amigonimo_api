package invitesctrl

// InviteInfoResponse represents information about an invite by code
type InviteInfoResponse struct {
	SecretFriendID string `json:"secretFriendId"`
	Name           string `json:"name"`
}
