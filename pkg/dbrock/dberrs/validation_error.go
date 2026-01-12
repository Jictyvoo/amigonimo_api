package dberrs

import (
	"fmt"
	"net/http"
)

// ErrDatabaseValidation represents a database validation error
// (e.g., invalid input data, missing required fields).
type ErrDatabaseValidation struct {
	baseErrorWrapper[ErrDatabaseValidation, *ErrDatabaseValidation]
	Field string
}

func NewErrDatabaseValidation(field string, err error) *ErrDatabaseValidation {
	newErr := &ErrDatabaseValidation{
		Field: field,
	}
	newErr.Err = err
	return newErr
}

func (e *ErrDatabaseValidation) Error() string {
	mainMessage := "database validation error"
	if e.Field != "" {
		mainMessage += fmt.Sprintf(", field '%s'", e.Field)
	}
	if e.Err != nil {
		return mainMessage + ": " + e.Err.Error()
	}
	return mainMessage
}

func (e *ErrDatabaseValidation) StatusCode() int {
	return http.StatusBadRequest
}
