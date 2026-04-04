//go:build integration

package wishlistrepo_test

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
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/wishlistrepo"
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

func genWishlistRepo(t testing.TB) *wishlistrepo.RepoMySQL {
	t.Helper()
	db := repoFactory.NewDB(t)
	return wishlistrepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(db))
}

type testEvent struct {
	user        entities.User
	sf          entities.SecretFriend
	participant entities.Participant
}

func seedMinimalEvent(t testing.TB, wishRepo *wishlistrepo.RepoMySQL, suffix string) testEvent {
	t.Helper()
	base := mysqlrepo.NewRepoMySQL(wishRepo.Connection())

	userID, err := entities.NewHexID()
	require.NoError(t, err)
	user := entities.User{
		ID: userID,
		UserBasic: authvalues.UserBasic{
			Username: "wish_user_" + suffix,
			Email:    "wish_" + suffix + "@example.com",
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
