package autherrs

import "net/http"

var (
	ErrWrongPassword = newError(
		"auth_wrong_password",
		http.StatusNotAcceptable,
		"provided password don't match with password for this user",
		nil,
		nil,
	)
	ErrUserNotFound = newError(
		"auth_user_not_found",
		http.StatusNotAcceptable,
		"user not found, make sure you have a valid session",
		nil,
		nil,
	)
	ErrEmailInUse = newError(
		"auth_email_in_use",
		http.StatusNotAcceptable,
		"this email is already in use and cannot be changed",
		nil,
		nil,
	)
	ErrUsernameInUse = newError(
		"auth_username_in_use",
		http.StatusNotAcceptable,
		"this username is already in use and cannot be changed",
		nil,
		nil,
	)
)

func NewErrChangeEmailLookup(err error) *Error {
	return newError(
		"auth_change_email_lookup_failed",
		http.StatusInternalServerError,
		"failed to validate email uniqueness",
		err,
		nil,
	)
}

func NewErrChangeUsernameLookup(err error) *Error {
	return newError(
		"auth_change_username_lookup_failed",
		http.StatusInternalServerError,
		"failed to validate username uniqueness",
		err,
		nil,
	)
}

func NewErrUpdatePassword(err error) *Error {
	return newError(
		"auth_update_password_failed",
		http.StatusInternalServerError,
		"internal error: cannot update the password",
		err,
		nil,
	)
}

func NewErrUpdateUsername(err error) *Error {
	return newError(
		"auth_update_username_failed",
		http.StatusInternalServerError,
		"internal error: cannot update the username",
		err,
		nil,
	)
}

func NewErrUpdateEmail(err error) *Error {
	return newError(
		"auth_update_email_failed",
		http.StatusInternalServerError,
		"internal error: cannot update the email",
		err,
		nil,
	)
}
