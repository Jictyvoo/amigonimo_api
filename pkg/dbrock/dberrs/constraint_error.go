package dberrs

import (
	"fmt"
	"net/http"
)

// ErrDatabaseConstraint represents a database constraint violation error
// (e.g., unique constraint, foreign key constraint)
type ErrDatabaseConstraint struct {
	baseErrorWrapper[ErrDatabaseConstraint, *ErrDatabaseConstraint]
	Constraint string
}

func NewErrDatabaseConstraint(constraint string, err error) (newErr *ErrDatabaseConstraint) {
	newErr = &ErrDatabaseConstraint{
		Constraint: constraint,
	}
	newErr.Err = err
	return newErr
}

func (e *ErrDatabaseConstraint) Error() string {
	mainMessage := "database constraint violation"
	if e.Err != nil {
		if e.Constraint != "" {
			mainMessage += fmt.Sprintf(", constraint '%s'", e.Constraint)
		}
		return mainMessage + ": " + e.Err.Error()
	}
	return mainMessage
}

func (e *ErrDatabaseConstraint) StatusCode() int {
	return http.StatusConflict
}
