//go:build integration

package denylistrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
)

func TestAddDenyListEntry_ByParticipantID(t *testing.T) {
	repo := genDenylistRepo(t)
	ev := seedMinimalEvent(t, repo, "add_by_pid")
	deniedID := seedExtraUser(t, repo, "add_by_pid")

	ref := denylist.ParticipantRef{ParticipantID: ev.participant.ID}
	entry, err := repo.AddDenyListEntry(ref, deniedID)
	require.NoError(t, err)
	assert.False(t, entry.ID.IsEmpty())
	assert.Equal(t, deniedID, entry.DeniedUserID)
}

func TestAddDenyListEntry_ByUserAndSF(t *testing.T) {
	repo := genDenylistRepo(t)
	ev := seedMinimalEvent(t, repo, "add_by_user_sf")
	deniedID := seedExtraUser(t, repo, "add_by_user_sf")

	ref := denylist.ParticipantRef{UserID: ev.user.ID, SecretFriendID: ev.sf.ID}
	entry, err := repo.AddDenyListEntry(ref, deniedID)
	require.NoError(t, err)
	assert.Equal(t, deniedID, entry.DeniedUserID)
}

func TestAddDenyListEntry_Duplicate(t *testing.T) {
	repo := genDenylistRepo(t)
	ev := seedMinimalEvent(t, repo, "add_dup")
	deniedID := seedExtraUser(t, repo, "add_dup")

	ref := denylist.ParticipantRef{ParticipantID: ev.participant.ID}
	_, err := repo.AddDenyListEntry(ref, deniedID)
	require.NoError(t, err)

	_, err = repo.AddDenyListEntry(ref, deniedID)
	assert.Error(t, err)
}
