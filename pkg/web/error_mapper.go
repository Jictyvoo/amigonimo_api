package web

import (
	"errors"
	"net/http"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if httpErr, ok := errors.AsType[*fuego.HTTPError](err); ok {
		return httpErr
	}

	if appErr, ok := errors.AsType[apperr.Contract](err); ok {
		return &fuego.HTTPError{
			Err:    appErr.Unwrap(),
			Type:   appErr.Code(),
			Title:  http.StatusText(appErr.StatusCode()),
			Status: appErr.StatusCode(),
			Detail: appErr.DetailMsg(),
		}
	}

	return newHTTPError(
		http.StatusInternalServerError,
		http.StatusText(http.StatusInternalServerError),
		err,
	)
}

func newHTTPError(statusCode int, detail string, err error) *fuego.HTTPError {
	return &fuego.HTTPError{
		Err:    err,
		Title:  http.StatusText(statusCode),
		Status: statusCode,
		Detail: detail,
	}
}
