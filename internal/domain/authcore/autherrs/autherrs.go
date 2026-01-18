package autherrs

import (
	"errors"
	"net/http"
)

type (
	errorInterface[T any] interface {
		*T
		error
	}
	baseErrorWrapper[Self interface{ reason() string }, SelfPtr errorInterface[Self]] struct {
		Err error
	}
)

func (e baseErrorWrapper[Self, SelfPtr]) Is(target error) bool {
	var asPointer SelfPtr
	if errors.As(target, &asPointer) {
		return true
	}

	_, ok := target.(Self) // Check the raw value directly
	return ok
}

func (e baseErrorWrapper[Self, SelfPtr]) Unwrap() error {
	return e.Err
}

func (e baseErrorWrapper[Self, SelfPtr]) StatusCode() int {
	return http.StatusInternalServerError
}

func (e baseErrorWrapper[Self, SelfPtr]) Error() string {
	var tempErr Self
	reason := tempErr.reason()
	if e.Err != nil {
		reason += ": " + e.Err.Error()
	}

	return reason
}

func (e baseErrorWrapper[Self, SelfPtr]) DetailMsg() string {
	var tempErr Self
	return tempErr.reason()
}

// Exported error instances for backward compatibility.
var (
	ErrEmailOrUsernameUsed  = &errEmailOrUsernameUsed{}
	ErrEmailUsed            = &errEmailUsed{}
	ErrUsernameUsed         = &errUsernameUsed{}
	ErrPasswordEncryption   = &errPasswordEncryption{}
	ErrUpdatePassword       = &errUpdatePassword{}
	ErrUpdateUsername       = &errUpdateUsername{}
	ErrWrongPassword        = &errWrongPassword{}
	ErrVerificationCode     = &errVerificationCode{}
	ErrGenRecoveryCode      = &errGenRecoveryCode{}
	ErrUserEmailNotFound    = &errUserEmailNotFound{}
	ErrUserRecoveryNotFound = &errUserRecoveryNotFound{}
	ErrInvalidAuthToken     = &errInvalidAuthToken{}
	ErrUpdateAuthToken      = &errUpdateAuthToken{}
)
