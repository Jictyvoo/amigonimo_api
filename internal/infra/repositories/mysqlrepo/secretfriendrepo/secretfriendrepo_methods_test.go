//go:build integration

package secretfriendrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func newSF(t testing.TB, ownerID entities.HexID, suffix string) *entities.SecretFriend {
	t.Helper()
	return &entities.SecretFriend{
		Name:            "sf-" + suffix,
		InviteCode:      "code-" + suffix,
		Status:          entities.StatusDraft,
		OwnerID:         ownerID,
		MaxDenyListSize: 5,
	}
}

func TestCreateSecretFriend_Success(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "create")

	sf := newSF(t, owner.ID, "create")
	err := repo.CreateSecretFriend(sf)
	require.NoError(t, err)
	assert.False(t, sf.ID.IsEmpty(), "ID should be assigned")
}

func TestGetSecretFriendByID_Success(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "get_id")

	sf := newSF(t, owner.ID, "get_id")
	require.NoError(t, repo.CreateSecretFriend(sf))

	got, err := repo.GetSecretFriendByID(sf.ID)
	require.NoError(t, err)
	assert.Equal(t, sf.ID, got.ID)
	assert.Equal(t, sf.Name, got.Name)
	assert.Equal(t, owner.ID, got.OwnerID)
}

func TestGetSecretFriendByID_NotFound(t *testing.T) {
	repo := genSecretFriendRepo(t)

	unknownID, err := entities.NewHexID()
	require.NoError(t, err)

	_, err = repo.GetSecretFriendByID(unknownID)
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestGetSecretFriendByInviteCode_Success(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "get_code")

	sf := newSF(t, owner.ID, "get_code")
	require.NoError(t, repo.CreateSecretFriend(sf))

	got, err := repo.GetSecretFriendByInviteCode(sf.InviteCode)
	require.NoError(t, err)
	assert.Equal(t, sf.ID, got.ID)
}

func TestGetSecretFriendByInviteCode_NotFound(t *testing.T) {
	repo := genSecretFriendRepo(t)

	_, err := repo.GetSecretFriendByInviteCode("unknown-code")
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestUpdateSecretFriend_Success(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "update")

	sf := newSF(t, owner.ID, "update")
	require.NoError(t, repo.CreateSecretFriend(sf))

	sf.Name = "updated-name"
	sf.Location = "new-location"
	require.NoError(t, repo.UpdateSecretFriend(sf))

	got, err := repo.GetSecretFriendByID(sf.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated-name", got.Name)
	assert.Equal(t, "new-location", got.Location)
}

func TestListSecretFriends_Success(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "list")

	sf1 := newSF(t, owner.ID, "list1")
	sf2 := newSF(t, owner.ID, "list2")
	require.NoError(t, repo.CreateSecretFriend(sf1))
	require.NoError(t, repo.CreateSecretFriend(sf2))

	list, err := repo.ListSecretFriends(owner.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
}
