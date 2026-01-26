package testengine

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http/httptest"
)

type Engine struct {
	ts     *httptest.Server
	db     *sql.DB
	dbName string
}

func (e *Engine) BaseURL() string {
	if e.ts != nil {
		return e.ts.URL
	}
	return "http://localhost:8080"
}

func (e *Engine) DB() *sql.DB {
	return e.db
}

func (e *Engine) Teardown() error {
	if e.ts != nil {
		e.ts.Close()
	}
	if e.db != nil {
		_, execErr := e.db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", e.dbName))
		closeErr := e.db.Close()
		if err := errors.Join(execErr, closeErr); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Seed(seeds ...any) {
	for _, s := range seeds {
		if err := executeSeedStruct(e.db, s); err != nil {
			panic(fmt.Errorf("failed to seed data: %w", err))
		}
	}
}
