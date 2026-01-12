package mysqlrepo

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

// WrapError wraps a database error into an appropriate custom error type
// based on the error content and MySQL error codes
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}

	// Check for common SQL errors
	if errors.Is(err, sql.ErrNoRows) {
		return dberrs.NewErrDatabaseNotFound(context, "", err)
	}

	errStr := err.Error()

	// MySQL error code patterns
	// Connection errors (1045, 2002, 2003, 2006, 2013)
	if strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "1045") ||
		strings.Contains(errStr, "2002") ||
		strings.Contains(errStr, "2003") ||
		strings.Contains(errStr, "2006") ||
		strings.Contains(errStr, "2013") {
		return dberrs.NewErrDatabaseConnection(err)
	}

	// Constraint violations (1062 = duplicate entry, 1452 = foreign key, 1451 = cannot delete)
	if strings.Contains(errStr, "1062") || // Duplicate entry
		strings.Contains(errStr, "1452") || // Foreign key constraint
		strings.Contains(errStr, "1451") || // Cannot delete or update
		strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "constraint") ||
		strings.Contains(errStr, "foreign key") {
		constraint := extractConstraintName(errStr)
		return dberrs.NewErrDatabaseConstraint(constraint, err)
	}

	// Validation errors (1048 = column cannot be null, 1264 = out of range)
	if strings.Contains(errStr, "1048") || // Column cannot be null
		strings.Contains(errStr, "1264") || // Out of range
		strings.Contains(errStr, "cannot be null") ||
		strings.Contains(errStr, "invalid") {
		field := extractFieldName(errStr)
		return dberrs.NewErrDatabaseValidation(field, err)
	}

	// Transaction errors (1213 = deadlock, 1205 = lock wait timeout)
	if strings.Contains(errStr, "1213") || // Deadlock
		strings.Contains(errStr, "1205") || // Lock wait timeout
		strings.Contains(errStr, "transaction") ||
		strings.Contains(errStr, "deadlock") ||
		strings.Contains(errStr, "rollback") {
		operation := extractOperation(context)
		return dberrs.NewErrDatabaseTransaction(operation, err)
	}

	// Default to query error
	return dberrs.NewErrDatabaseQuery(context, err)
}

// extractConstraintName attempts to extract the constraint name from error message
func extractConstraintName(errStr string) string {
	// Try to find constraint name in error message
	// MySQL format: "Duplicate entry 'value' for key 'constraint_name'"
	if idx := strings.Index(errStr, "for key '"); idx != -1 {
		start := idx + 9
		if end := strings.Index(errStr[start:], "'"); end != -1 {
			return errStr[start : start+end]
		}
	}
	return ""
}

// extractFieldName attempts to extract the field/column name from error message
func extractFieldName(errStr string) string {
	// MySQL format: "Column 'column_name' cannot be null"
	if idx := strings.Index(errStr, "Column '"); idx != -1 {
		start := idx + 8
		if end := strings.Index(errStr[start:], "'"); end != -1 {
			return errStr[start : start+end]
		}
	}
	return ""
}

// extractOperation extracts operation name from context string
func extractOperation(context string) string {
	// Simple extraction - could be enhanced
	if context == "" {
		return "unknown"
	}
	return context
}
