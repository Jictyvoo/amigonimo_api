//go:build integration

package denylistrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
)

func TestRemoveDenyListEntry_Success(t *testing.T) {
	repo := genDenylistRepo(t)
	ev := seedMinimalEvent(t, repo, "remove")
	deniedID := seedExtraUser(t, repo, "remove")

	ref := denylist.ParticipantRef{ParticipantID: ev.participant.ID}
	_, err := repo.AddDenyListEntry(ref, deniedID)
	require.NoError(t, err)

	err = repo.RemoveDenyListEntry(ref, deniedID)
	require.NoError(t, err)

	entries, err := repo.GetDenyListByParticipant(ref)
	require.NoError(t, err)
	assert.Empty(t, entries)
}
