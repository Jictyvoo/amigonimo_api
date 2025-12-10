package autherrs

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found, make sure you have a valid session")
	ErrEmailInUse    = errors.New("this email is already in use and cannot be changed")
	ErrUsernameInUse = errors.New("this username is already in use and cannot be changed")
)
