package mysqlrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

var _ dbrock.Transactioner = (*RepoMySQL)(nil)

// RepoMySQL is the main type to be used to perform database operations.
// It provides wrappers for every required.
type RepoMySQL struct {
	conn    *sql.DB
	queries *dbgen.Queries
}

func NewRepoMySQL(db *sql.DB) RepoMySQL {
	return RepoMySQL{conn: db, queries: dbgen.New(db)}
}

func (r *RepoMySQL) Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func (r *RepoMySQL) Connection() *sql.DB {
	return r.conn
}

func (r *RepoMySQL) Queries() *dbgen.Queries {
	return r.queries
}

func (r *RepoMySQL) BeginTx(
	ctx context.Context, txOpts *sql.TxOptions,
) (dbrock.OnFinishFunc, error) {
	tx, err := r.conn.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}

	r.queries = r.queries.WithTx(tx)
	txf := transactionFinisher{tx}
	return txf.Finish, nil
}

type transactionFinisher struct{ *sql.Tx }

func (tx transactionFinisher) Finish(commit bool) error {
	if commit {
		return tx.Commit()
	}

	return tx.Rollback()
}
