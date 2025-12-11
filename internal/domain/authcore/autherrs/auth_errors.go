package autherrs

import "net/http"

type errEmailOrUsernameUsed struct {
	baseErrorWrapper[errEmailOrUsernameUsed, *errEmailOrUsernameUsed]
}

func (e errEmailOrUsernameUsed) reason() string {
	return "user already with provided email/username already exists"
}

func (e errEmailOrUsernameUsed) StatusCode() int {
	return http.StatusPreconditionFailed
}

// errEmailUsed represents an error when email is already in use
type errEmailUsed struct {
	baseErrorWrapper[errEmailUsed, *errEmailUsed]
}

func (e errEmailUsed) reason() string {
	return "email already in use, if don't remember the password, access password recovery"
}

func (e errEmailUsed) StatusCode() int {
	return http.StatusPreconditionFailed
}

// errUsernameUsed represents an error when username is already in use
type errUsernameUsed struct {
	baseErrorWrapper[errUsernameUsed, *errUsernameUsed]
}

func (e errUsernameUsed) reason() string {
	return "username already in use, if don't remember the password, access password recovery"
}

func (e errUsernameUsed) StatusCode() int {
	return http.StatusPreconditionFailed
}

// errPasswordEncryption represents an error during password encryption
type errPasswordEncryption struct {
	baseErrorWrapper[errPasswordEncryption, *errPasswordEncryption]
}

func (e errPasswordEncryption) reason() string {
	return "password encryption error"
}

// errUpdatePassword represents an error when password update fails
type errUpdatePassword struct {
	baseErrorWrapper[errUpdatePassword, *errUpdatePassword]
}

func (e errUpdatePassword) reason() string {
	return "internal error: cannot update the password"
}

// errUpdateUsername represents an error when username update fails
type errUpdateUsername struct {
	baseErrorWrapper[errUpdateUsername, *errUpdateUsername]
}

func (e errUpdateUsername) reason() string {
	return "internal error: cannot update the username"
}

// errUserCreation represents an error when user creation fails
type errUserCreation struct {
	baseErrorWrapper[errUserCreation, *errUserCreation]
}

func (e errUserCreation) reason() string {
	return "user can't be created"
}

// errWrongPassword represents an error when password is incorrect
type errWrongPassword struct {
	baseErrorWrapper[errWrongPassword, *errWrongPassword]
}

func (e errWrongPassword) reason() string {
	return "provided password don't match with password for this user"
}

func (e errWrongPassword) StatusCode() int {
	return http.StatusNotAcceptable
}

// errUserEmailNotFound represents an error when user is not found by email/username
type errUserEmailNotFound struct {
	baseErrorWrapper[errUserEmailNotFound, *errUserEmailNotFound]
}

func (e errUserEmailNotFound) reason() string {
	return "not found user with given email/username and password combination"
}

func (e errUserEmailNotFound) StatusCode() int {
	return http.StatusNotAcceptable
}
