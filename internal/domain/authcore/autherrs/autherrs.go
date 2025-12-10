package autherrs

import (
	"errors"
)

var (
	ErrEmailUsed = errors.New(
		"email already in use, if don't remember the password, access password recovery",
	)
	ErrUsernameUsed = errors.New(
		"username already in use, if don't remember the password, access password recovery",
	)
	ErrPasswordEncryption = errors.New("password encryption error")
	ErrUpdatePassword     = errors.New("internal error: cannot update the password")
	ErrUpdateUsername     = errors.New("internal error: cannot update the username")
	ErrUserCreation       = errors.New("user can't be created")
	ErrWrongPassword      = errors.New(
		"provided password don't match with password for this user",
	)
	ErrVerificationCode = errors.New("cannot find any user with given code")
	ErrGenRecoveryCode  = errors.New(
		"internal error: cannot generate and save the recovery code",
	)
	ErrUserEmailNotFound = errors.New(
		"not found user with given email/username and password combination",
	)
	ErrUserRecoveryNotFound = errors.New("cannot find user with given email and recovery code")
	ErrInvalidAuthToken     = errors.New("error token provided is invalid")
	ErrUpdateAuthToken      = errors.New("was not possible to update the user authentication token")
)
