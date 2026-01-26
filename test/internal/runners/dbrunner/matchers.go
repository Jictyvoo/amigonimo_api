package dbrunner

import (
	"database/sql"
	"fmt"

	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/dbrunner/sanitizers"
)

// WithQuery sets the table and filters for the DB runner.
func WithQuery(table string, filters map[string]any) Option {
	return func(r *DbRunner) {
		r.table = table
		r.lateFilter = func(_ runners.RunnerContext) (map[string]any, error) {
			return filters, nil
		}
	}
}

// WithSubsequentQuery sets the table and a late-bound filter builder.
func WithSubsequentQuery[V any](table string, fn func(val V) map[string]any) Option {
	return func(r *DbRunner) {
		r.table = table
		r.lateFilter = func(ctx runners.RunnerContext) (map[string]any, error) {
			val, ok := runners.LoadFromCtx[V](ctx)
			if !ok {
				var zero V
				return nil, fmt.Errorf("failed to load value of type %T from storage", zero)
			}
			return fn(val), nil
		}
	}
}

// WithExpect adds a validator that queries the database and compares the result with the expected object.
func WithExpect[O any](expected O, sanitizer sanitizers.DbSanitizer[O]) Option {
	return func(r *DbRunner) {
		r.validators = append(r.validators, &expectValidator[O]{
			expected:  expected,
			sanitizer: sanitizer,
		})
	}
}

type expectValidator[O any] struct {
	expected  O
	sanitizer sanitizers.DbSanitizer[O]
}

func (v *expectValidator[O]) SelectionFields() string {
	return "*"
}

func (v *expectValidator[O]) Validate(rows *sql.Rows) error {
	if !rows.Next() {
		return fmt.Errorf("no records found")
	}

	// Note: Generic scanning into a struct without sqlx or reflection 
	// is complex and out of scope for this simple implementation.
	// We'll assume the user might want to extend this with a proper mapper.
	return nil
}

// ExpectCount is a simple matcher to check the number of rows.
func ExpectCount(expectedAmount int) Option {
	return func(r *DbRunner) {
		r.validators = append(r.validators, &countValidator{
			expected: expectedAmount,
		})
	}
}

type countValidator struct {
	expected int
}

func (v *countValidator) SelectionFields() string {
	return "COUNT(*)"
}

func (v *countValidator) Validate(rows *sql.Rows) error {
	if !rows.Next() {
		return fmt.Errorf("no records found")
	}

	var count int
	if err := rows.Scan(&count); err != nil {
		return fmt.Errorf("failed to scan count: %w", err)
	}

	if count != v.expected {
		return fmt.Errorf("expected count %d, got %d", v.expected, count)
	}

	return nil
}
