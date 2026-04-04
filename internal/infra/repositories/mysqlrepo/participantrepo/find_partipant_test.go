//go:build integration

package participantrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestGetParticipant_Success(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "get_ok")

	added, err := repo.AddParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)

	got, err := repo.GetParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)
	assert.Equal(t, added.ID, got.ID)
}

func TestGetParticipant_NotFound(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "get_nf")

	_, err := repo.GetParticipant(ev.sf.ID, ev.user.ID)
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestGetParticipantByID_Success(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "get_by_id")

	added, err := repo.AddParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)

	got, err := repo.GetParticipantByID(added.ID)
	require.NoError(t, err)
	assert.Equal(t, added.ID, got.ID)
}

func TestSetParticipantReady_Success(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "set_ready")

	_, err := repo.AddParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)

	err = repo.SetParticipantReady(ev.sf.ID, ev.user.ID, true)
	assert.NoError(t, err)
}

func TestRemoveParticipant_Success(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "remove")

	_, err := repo.AddParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)

	err = repo.RemoveParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)

	_, err = repo.GetParticipant(ev.sf.ID, ev.user.ID)
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}
