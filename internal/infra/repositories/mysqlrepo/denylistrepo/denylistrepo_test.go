//go:build integration

package denylistrepo_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authuserepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/denylistrepo"
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

// seedExtraUser creates a user in the same DB as the given deny repo.
// Used when a test needs a real user ID for the denied_user_id FK.
func seedExtraUser(t testing.TB, denyRepo *denylistrepo.RepoMySQL, suffix string) entities.HexID {
	t.Helper()
	id, err := entities.NewHexID()
	require.NoError(t, err)
	user := entities.User{
		ID: id,
		UserBasic: authvalues.UserBasic{
			Username: "extra_" + suffix,
			Email:    "extra_" + suffix + "@example.com",
			Password: "hashed",
		},
	}
	base := mysqlrepo.NewRepoMySQL(denyRepo.Connection())
	require.NoError(t, authuserepo.NewRepoMySQL(base).CreateUser(user, ""))
	return id
}

func genDenylistRepo(t testing.TB) *denylistrepo.RepoMySQL {
	t.Helper()
	db := repoFactory.NewDB(t)
	return denylistrepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(db))
}

type testEvent struct {
	user        entities.User
	sf          entities.SecretFriend
	participant entities.Participant
}

// seedMinimalEvent creates a user, secret friend, and participant sharing the same DB.
func seedMinimalEvent(t testing.TB, denyRepo *denylistrepo.RepoMySQL, suffix string) testEvent {
	t.Helper()
	base := mysqlrepo.NewRepoMySQL(denyRepo.Connection())

	userID, err := entities.NewHexID()
	require.NoError(t, err)
	user := entities.User{
		ID: userID,
		UserBasic: authvalues.UserBasic{
			Username: "deny_user_" + suffix,
			Email:    "deny_" + suffix + "@example.com",
			Password: "hashed",
		},
	}
	require.NoError(t, authuserepo.NewRepoMySQL(base).CreateUser(user, ""))

	sfRepo := secretfriendrepo.NewRepoMySQL(base)
	sf := &entities.SecretFriend{
		Name:            "test-sf-" + suffix,
		InviteCode:      "code-" + suffix,
		Status:          entities.StatusDraft,
		OwnerID:         userID,
		MaxDenyListSize: 5,
	}
	require.NoError(t, sfRepo.CreateSecretFriend(sf))

	pRepo := participantrepo.NewRepoMySQL(base)
	participant, err := pRepo.AddParticipant(sf.ID, userID)
	require.NoError(t, err)

	return testEvent{user: user, sf: *sf, participant: participant}
}
