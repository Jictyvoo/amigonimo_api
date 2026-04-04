//go:build integration

package authtokenrepo_test

import (
	"log"
	"os"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authtokenrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authuserepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/sqltest"
	"github.com/stretchr/testify/require"
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

func genAuthTokenRepo(t testing.TB) *authtokenrepo.RepoMySQL {
	t.Helper()
	db := repoFactory.NewDB(t)
	return authtokenrepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(db))
}

// seedUser creates a user in the same DB used by the token repo and returns the User.
func seedUser(t testing.TB, tokenRepo *authtokenrepo.RepoMySQL, suffix string) entities.User {
	t.Helper()
	id, err := entities.NewHexID()
	require.NoError(t, err)
	user := entities.User{
		ID: id,
		UserBasic: authvalues.UserBasic{
			Username: "tokenuser_" + suffix,
			Email:    "tokenuser_" + suffix + "@example.com",
			Password: "hashed",
		},
	}
	userRepo := authuserepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(tokenRepo.Connection()))
	require.NoError(t, userRepo.CreateUser(user, ""))
	return user
}
