package apperr

import (
	"errors"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func From(code, publicMessage string, err error) *Error {
	switch {
	case err == nil:
		return nil
	case isAppError(err):
		return fromAppError(code, publicMessage, err)
	default:
		return fromInfra(code, publicMessage, err)
	}
}

func fromAppError(code, publicMessage string, err error) *Error {
	existing, _ := errors.AsType[Contract](err)

	if code == "" && publicMessage == "" {
		if concrete, ok := errors.AsType[*Error](err); ok {
			return concrete
		}
	}

	if code == "" {
		code = existing.Code()
	}
	if publicMessage == "" {
		publicMessage = existing.Error()
	}

	internal := existing.Unwrap()
	if internal == nil {
		internal = err
	}

	return New(code, existing.StatusCode(), publicMessage, internal, existing.Metadata())
}

func fromInfra(code, publicMessage string, err error) *Error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, &dberrs.ErrDatabaseNotFound{}):
		return NotFound(code, publicMessage, err)
	case errors.Is(err, &dberrs.ErrDatabaseConstraint{}):
		return Conflict(code, publicMessage, err)
	case errors.Is(err, &dberrs.ErrDatabaseValidation{}):
		return Invalid(code, publicMessage, err)
	case errors.Is(err, &dberrs.ErrDatabaseConnection{}):
		return ServiceUnavailable(code, publicMessage, err)
	default:
		return InternalError(code, publicMessage, err)
	}
}

func isAppError(err error) bool {
	_, ok := errors.AsType[Contract](err)
	return ok
}
