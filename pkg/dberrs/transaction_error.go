package dberrs

import (
	"fmt"
)

// ErrDatabaseTransaction represents a database transaction error
type ErrDatabaseTransaction struct {
	baseErrorWrapper[ErrDatabaseTransaction, *ErrDatabaseTransaction]
	Operation string
}

func NewErrDatabaseTransaction(
	operation string, err error,
) *ErrDatabaseTransaction {
	newErr := &ErrDatabaseTransaction{
		Operation: operation,
	}
	newErr.Err = err
	return newErr
}

func (e *ErrDatabaseTransaction) Error() string {
	mainMessage := "database transaction error"
	if e.Operation != "" {
		mainMessage += fmt.Sprintf("operation '%s'", e.Operation)
	}
	if e.Err != nil {
		return mainMessage + ": " + e.Err.Error()
	}
	return mainMessage
}
