package autherrs

import "net/http"

var (
	ErrVerificationCode = newError(
		"auth_verification_code_invalid",
		http.StatusPreconditionFailed,
		"cannot find any user with given code",
		nil,
		nil,
	)
	ErrGenRecoveryCode = newError(
		"auth_generate_recovery_code_failed",
		http.StatusInternalServerError,
		"internal error: cannot generate and save the recovery code",
		nil,
		nil,
	)
	ErrUserEmailNotFound = newError(
		"auth_user_email_not_found",
		http.StatusNotAcceptable,
		"user not found",
		nil,
		nil,
	)
	ErrUserRecoveryNotFound = newError(
		"auth_user_recovery_not_found",
		http.StatusPreconditionFailed,
		"cannot find user with given email and recovery code",
		nil,
		nil,
	)
	ErrInvalidAuthToken = newError(
		"auth_invalid_token",
		http.StatusPreconditionFailed,
		"error token provided is invalid",
		nil,
		nil,
	)
	ErrUpdateAuthToken = newError(
		"auth_update_token_failed",
		http.StatusInternalServerError,
		"was not possible to update the user authentication token",
		nil,
		nil,
	)
)

func NewErrTokenLookup(err error) *Error {
	return newError(
		"auth_token_lookup_failed",
		http.StatusInternalServerError,
		"failed to load authentication token",
		err,
		nil,
	)
}

func NewErrTokenRegenerate(err error) *Error {
	return newError(
		"auth_token_regenerate_failed",
		http.StatusInternalServerError,
		"failed to regenerate authentication token",
		err,
		nil,
	)
}

func NewErrUpdateAuthToken(err error) *Error {
	return newError(
		"auth_update_token_failed",
		http.StatusInternalServerError,
		"was not possible to update the user authentication token",
		err,
		nil,
	)
}

func NewErrGenRecoveryCode(err error) *Error {
	return newError(
		"auth_generate_recovery_code_failed",
		http.StatusInternalServerError,
		"internal error: cannot generate and save the recovery code",
		err,
		nil,
	)
}

func NewErrSetVerification(err error) *Error {
	return newError(
		"auth_set_verification_failed",
		http.StatusInternalServerError,
		"cannot update the verification state",
		err,
		nil,
	)
}
