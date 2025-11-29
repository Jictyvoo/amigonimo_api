package mysqlrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

type RepoMySQL struct {
	conn    *sql.DB
	queries *dbgen.Queries
}

func NewRepoMySQL(db *sql.DB) RepoMySQL {
	return RepoMySQL{conn: db, queries: dbgen.New(db)}
}

func (r RepoMySQL) Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
