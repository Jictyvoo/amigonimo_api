//go:build integration

package authtokenrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func TestUpsertAuthToken_CreateAndUpdate(t *testing.T) {
	repo := genAuthTokenRepo(t)
	user := seedUser(t, repo, "upsert")

	token := &entities.AuthenticationToken{
		User: user,
	}
	require.NoError(t, token.Regenerate(0))
	token.User = user
	require.NoError(t, repo.UpsertAuthToken(token))
	assert.False(t, token.ID.IsEmpty(), "ID should be assigned after first upsert")

	firstToken := token.AuthToken

	// Second upsert should update (same user, different token string)
	require.NoError(t, token.Regenerate(0))
	token.User = user
	require.NoError(t, repo.UpsertAuthToken(token))
	assert.NotEqual(t, firstToken, token.AuthToken)
}
