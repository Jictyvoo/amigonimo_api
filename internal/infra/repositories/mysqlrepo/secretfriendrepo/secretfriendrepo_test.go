//go:build integration

package secretfriendrepo_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authuserepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/sqltest"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/participantrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/secretfriendrepo"
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

func genSecretFriendRepo(t testing.TB) *secretfriendrepo.RepoMySQL {
	t.Helper()
	db := repoFactory.NewDB(t)
	return secretfriendrepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(db))
}

func seedOwner(t testing.TB, sfRepo *secretfriendrepo.RepoMySQL, suffix string) entities.User {
	t.Helper()
	id, err := entities.NewHexID()
	require.NoError(t, err)
	user := entities.User{
		ID: id,
		UserBasic: authvalues.UserBasic{
			Username: "sf_owner_" + suffix,
			Email:    "sf_owner_" + suffix + "@example.com",
			Password: "hashed",
		},
	}
	base := mysqlrepo.NewRepoMySQL(sfRepo.Connection())
	require.NoError(t, authuserepo.NewRepoMySQL(base).CreateUser(user, ""))
	return user
}

func seedParticipant(
	t testing.TB,
	sfRepo *secretfriendrepo.RepoMySQL,
	sfID, userID entities.HexID,
) entities.Participant {
	t.Helper()
	base := mysqlrepo.NewRepoMySQL(sfRepo.Connection())
	p, err := participantrepo.NewRepoMySQL(base).AddParticipant(sfID, userID)
	require.NoError(t, err)
	return p
}
