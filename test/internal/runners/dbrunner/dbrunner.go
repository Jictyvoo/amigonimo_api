package dbrunner

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
)

type dbValidator interface {
	SelectionFields() string
	Validate(rows *sql.Rows) error
}

// DbRunner is a test runner for database assertions.
type DbRunner struct {
	db         *sql.DB
	table      string
	lateFilter func(runners.RunnerContext) (map[string]any, error)
	validators []dbValidator
}

// Option is a functional option for configuring the DbRunner.
type Option func(*DbRunner)

// NewDbRunner creates a new DbRunner with the given options.
func NewDbRunner(db *sql.DB, opts ...Option) *DbRunner {
	r := &DbRunner{
		db: db,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *DbRunner) Run(rCtx runners.RunnerContext) error {
	var filters map[string]any
	if r.lateFilter != nil {
		var err error
		filters, err = r.lateFilter(rCtx)
		if err != nil {
			rCtx.Fatalf("Failed to build late filters: %v", err)
		}
	}

	var (
		where = make([]string, 0, len(filters))
		args  = make([]any, 0, len(filters))
	)
	for filterKey, filterValue := range filters {
		where = append(where, fmt.Sprintf("%s = ?", filterKey))
		args = append(args, filterValue)
	}
	for _, v := range r.validators {
		query := fmt.Sprintf(
			"SELECT %s FROM %s WHERE %s",
			v.SelectionFields(), r.table, strings.Join(where, " AND "),
		)

		rows, err := r.db.Query(query, args...)
		if err != nil {
			rCtx.Fatalf("Failed to execute query %q: %v", query, err)
			return err
		}

		if err = v.Validate(rows); err != nil {
			rCtx.Fatalf("Database validation failed: %v", err)
		}
		_ = rows.Close()
	}
	return nil
}
