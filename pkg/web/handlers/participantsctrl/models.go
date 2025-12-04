package participantsctrl

// ConfirmParticipationRequest represents the request to confirm or leave participation
type ConfirmParticipationRequest struct {
	Confirm bool `json:"confirm" validate:"required"` // true to join, false to leave
}

// ConfirmParticipationResponse represents the response when confirming participation
type ConfirmParticipationResponse struct {
	Success       bool   `json:"success"`
	ParticipantID string `json:"participantId"`
}

// ParticipantResponse represents a participant in the list
type ParticipantResponse struct {
	ParticipantID string `json:"participantId"`
	UserID        string `json:"userId"`
	Fullname      string `json:"fullname"`
}
