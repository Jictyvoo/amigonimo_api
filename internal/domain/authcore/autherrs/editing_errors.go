package autherrs

import "net/http"

// errUserNotFound represents an error when user is not found.
type errUserNotFound struct {
	baseErrorWrapper[errUserNotFound, *errUserNotFound]
}

func (e errUserNotFound) StatusCode() int {
	return http.StatusNotAcceptable
}

func (e errUserNotFound) reason() string {
	return "user not found, make sure you have a valid session"
}

// errEmailInUse represents an error when email is already in use.
type errEmailInUse struct {
	baseErrorWrapper[errEmailInUse, *errEmailInUse]
}

func (e errEmailInUse) StatusCode() int {
	return http.StatusNotAcceptable
}

func (e errEmailInUse) reason() string {
	return "this email is already in use and cannot be changed"
}

// errUsernameInUse represents an error when username is already in use.
type errUsernameInUse struct {
	baseErrorWrapper[errUsernameInUse, *errUsernameInUse]
}

func (e errUsernameInUse) StatusCode() int {
	return http.StatusNotAcceptable
}

func (e errUsernameInUse) reason() string {
	return "this username is already in use and cannot be changed"
}

// Exported error instances for backward compatibility.
var (
	ErrUserNotFound  = &errUserNotFound{}
	ErrEmailInUse    = &errEmailInUse{}
	ErrUsernameInUse = &errUsernameInUse{}
)
