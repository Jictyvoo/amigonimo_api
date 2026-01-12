package dashboardctrl

// SecretFriendSummary represents a summary of a secret friend event.
type SecretFriendSummary struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Datetime          *string `json:"datetime,omitempty"`
	Location          string  `json:"location,omitempty"`
	Status            string  `json:"status"`
	ParticipantsCount int     `json:"participantsCount"`
}

// DashboardResponse represents the dashboard data for a user.
type DashboardResponse struct {
	ActiveCreated     []SecretFriendSummary `json:"activeCreated"`
	ActiveParticipant []SecretFriendSummary `json:"activeParticipant"`
}
