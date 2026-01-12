package dbrock

import (
	"context"
	"database/sql"
)

type OnFinishFunc func(commit bool) error

type Transactioner interface {
	BeginTx(ctx context.Context, txOpts *sql.TxOptions) (OnFinishFunc, error)
}
