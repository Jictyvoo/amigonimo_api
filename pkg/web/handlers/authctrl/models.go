package authctrl

// LoginRequest represents the request to login
type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// LoginResponse represents the response when logging in
type LoginResponse struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

// RegisterRequest represents the request to register
type RegisterRequest struct {
	FullName   string `json:"fullname"             validate:"required"`
	Email      string `json:"email"                validate:"required,email"`
	InviteCode string `json:"inviteCode,omitempty"`
}

// RegisterResponse represents the response when registering
type RegisterResponse struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}
