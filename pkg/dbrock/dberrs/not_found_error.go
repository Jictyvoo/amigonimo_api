package dberrs

import (
	"fmt"
	"net/http"
)

// ErrDatabaseNotFound represents a "not found" error (e.g., sql.ErrNoRows).
type ErrDatabaseNotFound struct {
	baseErrorWrapper[ErrDatabaseNotFound, *ErrDatabaseNotFound]

	Resource   string
	Identifier string
}

func NewErrDatabaseNotFound(
	resource, identifier string, err error,
) *ErrDatabaseNotFound {
	newErr := &ErrDatabaseNotFound{
		Resource:   resource,
		Identifier: identifier,
	}
	newErr.Err = err
	return newErr
}

func (e *ErrDatabaseNotFound) Error() string {
	mainMessage := "resource not found"
	if e.Resource != "" && e.Identifier != "" {
		mainMessage = fmt.Sprintf(
			"%s with identifier '%s' not found",
			e.Resource, e.Identifier,
		)
	}
	if e.Err != nil {
		return mainMessage + ": " + e.Err.Error()
	}
	return mainMessage
}

func (e *ErrDatabaseNotFound) StatusCode() int {
	return http.StatusNotFound
}
