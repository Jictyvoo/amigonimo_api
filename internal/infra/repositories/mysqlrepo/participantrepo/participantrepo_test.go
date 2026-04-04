//go:build integration

package participantrepo_test

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

func genParticipantRepo(t testing.TB) *participantrepo.RepoMySQL {
	t.Helper()
	db := repoFactory.NewDB(t)
	return participantrepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(db))
}

type seedResult struct {
	user entities.User
	sf   entities.SecretFriend
}

func seedUserAndSF(t testing.TB, pRepo *participantrepo.RepoMySQL, suffix string) seedResult {
	t.Helper()
	base := mysqlrepo.NewRepoMySQL(pRepo.Connection())

	userID, err := entities.NewHexID()
	require.NoError(t, err)
	user := entities.User{
		ID: userID,
		UserBasic: authvalues.UserBasic{
			Username: "p_user_" + suffix,
			Email:    "p_" + suffix + "@example.com",
			Password: "hashed",
		},
	}
	require.NoError(t, authuserepo.NewRepoMySQL(base).CreateUser(user, ""))

	sfRepo := secretfriendrepo.NewRepoMySQL(base)
	sf := &entities.SecretFriend{
		Name:            "sf-" + suffix,
		InviteCode:      "inv-" + suffix,
		Status:          entities.StatusDraft,
		OwnerID:         userID,
		MaxDenyListSize: 3,
	}
	require.NoError(t, sfRepo.CreateSecretFriend(sf))

	return seedResult{user: user, sf: *sf}
}
