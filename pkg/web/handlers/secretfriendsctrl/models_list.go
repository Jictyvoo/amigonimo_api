package secretfriendsctrl

import "time"

type SecretFriendSummary struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Datetime          time.Time `json:"datetime,omitempty"`
	Location          string    `json:"location,omitempty"`
	Status            string    `json:"status"`
	ParticipantsCount uint8     `json:"participantsCount"`
}

type (
	EventsList struct {
		Created     []SecretFriendSummary `json:"created"`
		Participant []SecretFriendSummary `json:"participant"`
	}
	DashboardResponse struct {
		Active   EventsList `json:"active"`
		Inactive EventsList `json:"inactive"`
	}
)
