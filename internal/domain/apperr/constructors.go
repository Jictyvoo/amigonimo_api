package apperr

import "net/http"

func Invalid(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusBadRequest, publicMessage, internal, nil)
}

func Unauthorized(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusUnauthorized, publicMessage, internal, nil)
}

func Forbidden(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusForbidden, publicMessage, internal, nil)
}

func NotFound(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusNotFound, publicMessage, internal, nil)
}

func Conflict(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusConflict, publicMessage, internal, nil)
}

func InternalError(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusInternalServerError, publicMessage, internal, nil)
}

func ServiceUnavailable(code, publicMessage string, internal error) *Error {
	return New(code, http.StatusServiceUnavailable, publicMessage, internal, nil)
}
