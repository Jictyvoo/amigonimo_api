package entities

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginResponse struct {
	UserID HexID  `json:"-"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	FullName   string `json:"fullname"              validate:"required"`
	Email      string `json:"email"                 validate:"required,email"`
	InviteCode string `json:"invite_code,omitempty"`
}

type RegisterResponse = LoginResponse
