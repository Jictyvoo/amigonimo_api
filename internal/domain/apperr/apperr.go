package apperr

import (
	"maps"
	"net/http"
)

type Contract interface {
	error
	Code() string
	StatusCode() int
	DetailMsg() string
	Metadata() map[string]any
	Unwrap() error
}

var _ Contract = (*Error)(nil)

type Error struct {
	code          string
	statusCode    int
	publicMessage string
	internalErr   error
	metadata      map[string]any
}

func New(
	code string,
	statusCode int,
	publicMessage string,
	internalErr error,
	metadata map[string]any,
) *Error {
	if statusCode <= 0 {
		statusCode = http.StatusInternalServerError
	}
	if publicMessage == "" {
		publicMessage = http.StatusText(statusCode)
	}

	return &Error{
		code:          code,
		statusCode:    statusCode,
		publicMessage: publicMessage,
		internalErr:   internalErr,
		metadata:      maps.Clone(metadata),
	}
}

func (e *Error) Error() string {
	return e.publicMessage
}

func (e *Error) Unwrap() error {
	return e.internalErr
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) StatusCode() int {
	return e.statusCode
}

func (e *Error) DetailMsg() string {
	return e.publicMessage
}

func (e *Error) Metadata() map[string]any {
	return maps.Clone(e.metadata)
}
