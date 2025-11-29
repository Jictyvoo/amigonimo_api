package dberrs

import (
	"errors"
	"net/http"
)

type (
	errorInterface[T any] interface {
		*T
		error
	}
	baseErrorWrapper[Self any, SelfPtr errorInterface[Self]] struct {
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

func (e baseErrorWrapper[Self, SelfPtr]) DetailMsg() string {
	if e.Err == nil {
		return "unknown error"
	}
	return e.Err.Error()
}
