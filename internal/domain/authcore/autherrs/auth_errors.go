package autherrs

import "net/http"

var (
	ErrEmailOrUsernameUsed = newError(
		"auth_email_or_username_used",
		http.StatusPreconditionFailed,
		"user already with provided email/username already exists",
		nil,
		nil,
	)
	ErrEmailUsed = newError(
		"auth_email_used",
		http.StatusPreconditionFailed,
		"email already in use, if don't remember the password, access password recovery",
		nil,
		nil,
	)
	ErrUsernameUsed = newError(
		"auth_username_used",
		http.StatusPreconditionFailed,
		"username already in use, if don't remember the password, access password recovery",
		nil,
		nil,
	)
	ErrInvalidCredentials = newError(
		"auth_invalid_credentials",
		http.StatusNotAcceptable,
		"not found user with given email/username and password combination",
		nil,
		nil,
	)
)

func NewErrSignUpLookup(err error) *Error {
	return newError(
		"auth_signup_lookup_failed",
		http.StatusInternalServerError,
		"failed to validate user uniqueness",
		err,
		nil,
	)
}

func NewErrLogin(err error) *Error {
	return newError(
		"auth_login_failed",
		http.StatusInternalServerError,
		"failed to log in",
		err,
		nil,
	)
}

func NewErrRecoveryLookup(err error) *Error {
	return newError(
		"auth_recovery_lookup_failed",
		http.StatusInternalServerError,
		"failed to look up recovery contact",
		err,
		nil,
	)
}

func NewErrPasswordEncryption(err error) *Error {
	return newError(
		"auth_password_encryption_failed",
		http.StatusInternalServerError,
		"password encryption error",
		err,
		nil,
	)
}

func NewErrUserCreation(err error) *Error {
	return newError(
		"auth_user_creation_failed",
		http.StatusInternalServerError,
		"user can't be created",
		err,
		nil,
	)
}
