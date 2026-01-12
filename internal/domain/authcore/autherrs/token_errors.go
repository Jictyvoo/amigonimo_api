package autherrs

import "net/http"

// errUserRecoveryNotFound represents an error when user recovery is not found.
type errUserRecoveryNotFound struct {
	baseErrorWrapper[errUserRecoveryNotFound, *errUserRecoveryNotFound]
}

func (e errUserRecoveryNotFound) StatusCode() int {
	return http.StatusPreconditionFailed
}

func (e errUserRecoveryNotFound) reason() string {
	return "cannot find user with given email and recovery code"
}

// errInvalidAuthToken represents an error when authentication token is invalid.
type errInvalidAuthToken struct {
	baseErrorWrapper[errInvalidAuthToken, *errInvalidAuthToken]
}

func (e errInvalidAuthToken) StatusCode() int {
	return http.StatusPreconditionFailed
}

func (e errInvalidAuthToken) reason() string {
	return "error token provided is invalid"
}

// errUpdateAuthToken represents an error when auth token update fails.
type errUpdateAuthToken struct {
	baseErrorWrapper[errUpdateAuthToken, *errUpdateAuthToken]
}

func (e errUpdateAuthToken) reason() string {
	return "was not possible to update the user authentication token"
}

// errGenRecoveryCode represents an error when recovery code generation fails.
type errGenRecoveryCode struct {
	baseErrorWrapper[errGenRecoveryCode, *errGenRecoveryCode]
}

func (e errGenRecoveryCode) reason() string {
	return "internal error: cannot generate and save the recovery code"
}

// errVerificationCode represents an error when verification code is invalid.
type errVerificationCode struct {
	baseErrorWrapper[errVerificationCode, *errVerificationCode]
}

func (e errVerificationCode) StatusCode() int {
	return http.StatusPreconditionFailed
}

func (e errVerificationCode) reason() string {
	return "cannot find any user with given code"
}
