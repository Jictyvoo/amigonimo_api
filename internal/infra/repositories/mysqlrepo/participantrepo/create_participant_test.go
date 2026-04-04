//go:build integration

package participantrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddParticipant_Success(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "add")

	p, err := repo.AddParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)
	assert.False(t, p.ID.IsEmpty())
	assert.Equal(t, ev.sf.ID, p.SecretFriendID)
	assert.Equal(t, ev.user.ID, p.RelatedUser.ID)
}
