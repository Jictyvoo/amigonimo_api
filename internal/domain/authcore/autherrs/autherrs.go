package autherrs

import (
	"maps"
	"net/http"
)

type Error struct {
	code          string
	statusCode    int
	publicMessage string
	internalErr   error
	metadata      map[string]any
}

func newError(
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

func (e *Error) Is(target error) bool {
	other, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.code != "" && e.code == other.code
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
