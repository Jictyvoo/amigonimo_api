package controllers

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

// FormEditEmail represents the request to edit user email
type FormEditEmail struct {
	CurrentPassword string `json:"current_password"`
	NewEmail        string `json:"new_email"`
}

// FormEditPassword represents the request to edit user password
type FormEditPassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// FormEditUsername represents the request to edit user username
type FormEditUsername struct {
	CurrentPassword string `json:"current_password"`
	NewUsername     string `json:"new_username"`
}

// FormRecoveryCode represents the request with recovery code
type FormRecoveryCode struct {
	Email        string `json:"email"`
	RecoveryCode string `json:"recovery_code"`
}

// FormUser represents the user form for signup/login
type FormUser struct {
	Username string `json:"username"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password"`
}

// FormResetPassword represents the request to reset password
type FormResetPassword struct {
	Email        string `json:"email"`
	RecoveryCode string `json:"recovery_code"`
	NewPassword  string `json:"new_password"`
}

// SuccessResponse represents a successful operation response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ForgotPasswordResponse represents the response for password recovery initiation
type ForgotPasswordResponse struct {
	SuccessResponse

	ObfuscatedEmail string `json:"obfuscatedEmail,omitempty"`
}
