package dberrs

import "net/http"

// ErrDatabaseConnection represents a database connection error.
type ErrDatabaseConnection struct {
	baseErrorWrapper[ErrDatabaseConnection, *ErrDatabaseConnection]
}

func NewErrDatabaseConnection(err error) (newErr *ErrDatabaseConnection) {
	newErr = new(ErrDatabaseConnection)
	newErr.Err = err
	return newErr
}

func (e *ErrDatabaseConnection) Error() string {
	mainMessage := "database connection failed"

	if e.Err != nil {
		mainMessage = mainMessage + ": " + e.Err.Error()
	}
	return mainMessage
}

func (e *ErrDatabaseConnection) StatusCode() int {
	return http.StatusServiceUnavailable
}
