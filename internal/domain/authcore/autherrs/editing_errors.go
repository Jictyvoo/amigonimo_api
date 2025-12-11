package autherrs

import "net/http"

// errUserNotFound represents an error when user is not found
type errUserNotFound struct {
	baseErrorWrapper[errUserNotFound, *errUserNotFound]
}

func (e errUserNotFound) reason() string {
	return "user not found, make sure you have a valid session"
}

func (e errUserNotFound) StatusCode() int {
	return http.StatusNotAcceptable
}

// errEmailInUse represents an error when email is already in use
type errEmailInUse struct {
	baseErrorWrapper[errEmailInUse, *errEmailInUse]
}

func (e errEmailInUse) reason() string {
	return "this email is already in use and cannot be changed"
}

func (e errEmailInUse) StatusCode() int {
	return http.StatusNotAcceptable
}

// errUsernameInUse represents an error when username is already in use
type errUsernameInUse struct {
	baseErrorWrapper[errUsernameInUse, *errUsernameInUse]
}

func (e errUsernameInUse) reason() string {
	return "this username is already in use and cannot be changed"
}

func (e errUsernameInUse) StatusCode() int {
	return http.StatusNotAcceptable
}

// Exported error instances for backward compatibility
var (
	ErrUserNotFound  = &errUserNotFound{}
	ErrEmailInUse    = &errEmailInUse{}
	ErrUsernameInUse = &errUsernameInUse{}
)
