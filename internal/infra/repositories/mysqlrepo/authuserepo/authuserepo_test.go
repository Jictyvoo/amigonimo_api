//go:build integration

package authuserepo_test

import (
	"log"
	"os"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authuserepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/sqltest"
)

var repoFactory *sqltest.Factory

func TestMain(m *testing.M) {
	var err error
	if repoFactory, err = sqltest.NewFactory(); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	if err = repoFactory.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func genAuthUserRepo(t testing.TB) *authuserepo.RepoMySQL {
	t.Helper()
	db := repoFactory.NewDB(t)
	return authuserepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(db))
}
