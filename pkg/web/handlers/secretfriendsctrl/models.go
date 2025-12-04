package secretfriendsctrl

import "time"

// CreateSecretFriendRequest represents the request to create a secret friend
type CreateSecretFriendRequest struct {
	Name            string     `json:"name"                      validate:"required"`
	Datetime        *time.Time `json:"datetime,omitempty"`
	Location        string     `json:"location,omitempty"`
	MaxDenyListSize *int       `json:"maxDenyListSize,omitempty"`
}

// CreateSecretFriendResponse represents the response when creating a secret friend
type CreateSecretFriendResponse struct {
	SecretFriendID string `json:"secretFriendId"`
	InviteCode     string `json:"inviteCode"`
	InviteLink     string `json:"inviteLink"`
}

// GetSecretFriendResponse represents the response when getting secret friend details
type GetSecretFriendResponse struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Datetime          *time.Time `json:"datetime,omitempty"`
	Location          string     `json:"location,omitempty"`
	OwnerID           string     `json:"ownerId"`
	ParticipantsCount int        `json:"participantsCount"`
	Status            string     `json:"status"` // draft | open | drawn | closed
}

// UpdateSecretFriendRequest represents the request to update a secret friend
type UpdateSecretFriendRequest struct {
	Name     *string    `json:"name,omitempty"`
	Datetime *time.Time `json:"datetime,omitempty"`
	Location *string    `json:"location,omitempty"`
}

// DrawSecretFriendResponse represents the response when drawing a secret friend
type DrawSecretFriendResponse struct {
	SecretFriendID string `json:"secretFriendId"`
	Status         string `json:"status"`
	ResultCount    int    `json:"resultCount"`
}
