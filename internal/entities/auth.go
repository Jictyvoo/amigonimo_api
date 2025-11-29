package entities

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type LoginResponse struct {
	UserID HexID  `json:"userId"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	Fullname   string `json:"fullname"             validate:"required"`
	Email      string `json:"email"                validate:"required,email"`
	InviteCode string `json:"inviteCode,omitempty"`
}

type RegisterResponse struct {
	UserID HexID  `json:"userId"`
	Token  string `json:"token"`
}

type EnterByCodeRequest struct {
	InviteCode string `json:"inviteCode" validate:"required"`
}

type EnterByCodeResponse struct {
	SecretFriendID HexID `json:"secretFriendId"`
	RequiresLogin  bool  `json:"requiresLogin"`
}
