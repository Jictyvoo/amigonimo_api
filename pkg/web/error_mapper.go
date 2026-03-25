package web

import (
	"errors"
	"net/http"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
)

type mappedError interface {
	error
	Code() string
	StatusCode() int
	DetailMsg() string
	Unwrap() error
}

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if httpErr, ok := errors.AsType[*fuego.HTTPError](err); ok {
		return httpErr
	}

	if appErr, ok := errors.AsType[apperr.Contract](err); ok {
		return mapDomainError(appErr)
	}

	if authErr, ok := errors.AsType[*autherrs.Error](err); ok {
		return mapDomainError(authErr)
	}

	return newHTTPError(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		err,
	)
}

func mapDomainError(err mappedError) *fuego.HTTPError {
	return &fuego.HTTPError{
		Err:    err.Unwrap(),
		Type:   err.Code(),
		Title:  http.StatusText(err.StatusCode()),
		Status: err.StatusCode(),
		Detail: err.DetailMsg(),
	}
}

func newHTTPError(statusCode int, detail string, err error) *fuego.HTTPError {
	return &fuego.HTTPError{
		Err:    err,
		Title:  http.StatusText(statusCode),
		Status: statusCode,
		Detail: detail,
	}
}
