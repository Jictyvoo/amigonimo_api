//go:build integration

package denylistrepo_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
)

func TestGetDenyListByParticipant_Empty(t *testing.T) {
	repo := genDenylistRepo(t)
	ev := seedMinimalEvent(t, repo, "list_empty")

	ref := denylist.ParticipantRef{ParticipantID: ev.participant.ID}
	entries, err := repo.GetDenyListByParticipant(ref)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

func TestGetDenyListByParticipant_WithEntries(t *testing.T) {
	repo := genDenylistRepo(t)
	ev := seedMinimalEvent(t, repo, "list_with")

	ref := denylist.ParticipantRef{ParticipantID: ev.participant.ID}

	for i := range 3 {
		deniedID := seedExtraUser(t, repo, fmt.Sprintf("list_with_%d", i))
		_, err := repo.AddDenyListEntry(ref, deniedID)
		require.NoError(t, err)
	}

	entries, err := repo.GetDenyListByParticipant(ref)
	require.NoError(t, err)
	assert.Len(t, entries, 3)
}
