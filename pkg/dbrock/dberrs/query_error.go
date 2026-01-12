package dberrs

import "fmt"

// ErrDatabaseQuery represents a database query execution error
type ErrDatabaseQuery struct {
	baseErrorWrapper[ErrDatabaseQuery, *ErrDatabaseQuery]
	Query string
}

func NewErrDatabaseQuery(query string, err error) *ErrDatabaseQuery {
	newErr := &ErrDatabaseQuery{
		Query: query,
	}
	newErr.Err = err
	return newErr
}

func (e *ErrDatabaseQuery) Error() string {
	mainMessage := "database query error"
	if e.Query != "" {
		mainMessage += fmt.Sprintf(", query '%s'", e.Query)
	}
	if e.Err != nil {
		return mainMessage + ": " + e.Err.Error()
	}
	return mainMessage
}
