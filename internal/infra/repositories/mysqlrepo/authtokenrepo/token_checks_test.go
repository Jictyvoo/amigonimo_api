//go:build integration

package authtokenrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestGetAuthenticationToken_NotFound_ReturnsEmpty(t *testing.T) {
	repo := genAuthTokenRepo(t)

	unknownID, err := entities.NewHexID()
	require.NoError(t, err)

	// GetAuthenticationToken returns empty token with no error when not found
	got, err := repo.GetAuthenticationToken(unknownID)
	require.NoError(t, err)
	assert.True(t, got.ID.IsEmpty())
}

func TestGetUserByAuthToken_Success(t *testing.T) {
	repo := genAuthTokenRepo(t)
	user := seedUser(t, repo, "byauthtoken")

	token := &entities.AuthenticationToken{User: user}
	require.NoError(t, token.Regenerate(0))
	token.User = user
	require.NoError(t, repo.UpsertAuthToken(token))

	got, err := repo.GetUserByAuthToken(token.AuthToken)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestGetUserByAuthToken_NotFound(t *testing.T) {
	repo := genAuthTokenRepo(t)

	_, err := repo.GetUserByAuthToken("nonexistent-token")
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestCheckAuthenticationByRefreshToken_Success(t *testing.T) {
	repo := genAuthTokenRepo(t)
	user := seedUser(t, repo, "byrefresh")

	token := &entities.AuthenticationToken{User: user}
	require.NoError(t, token.Regenerate(0))
	token.User = user
	require.NoError(t, repo.UpsertAuthToken(token))

	// Refresh token is stored as raw bytes (string(uuid[:])), not canonical UUID string
	refreshStr := string(token.RefreshToken.UUID[:])
	got, err := repo.CheckAuthenticationByRefreshToken(refreshStr)
	require.NoError(t, err)
	assert.Equal(t, token.AuthToken, got.AuthToken)
}

func TestCheckAuthenticationByRefreshToken_NotFound(t *testing.T) {
	repo := genAuthTokenRepo(t)

	_, err := repo.CheckAuthenticationByRefreshToken("wrong-refresh-token")
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}
