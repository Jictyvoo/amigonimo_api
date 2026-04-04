//go:build integration

package authuserepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestGetUserByEmail_NotFound(t *testing.T) {
	repo := genAuthUserRepo(t)
	_, err := repo.GetUserByEmail("nonexistent@example.com")
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestGetUserByEmailOrUsername_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "byemailorusername")
	require.NoError(t, repo.CreateUser(user, ""))

	got, err := repo.GetUserByEmailOrUsername(user.Email, user.Username)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestGetUserByUsername_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "byusername")
	require.NoError(t, repo.CreateUser(user, ""))

	got, err := repo.GetUserByUsername(user.Username)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestGetUserByVerificationCode_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "byverification")
	const token = "verification-code-abc"
	require.NoError(t, repo.CreateUser(user, token))

	got, err := repo.GetUserByVerificationCode(token)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestGetUserByVerificationCode_NotFound(t *testing.T) {
	repo := genAuthUserRepo(t)
	_, err := repo.GetUserByVerificationCode("wrong-code")
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}

func TestGetUserByRecovery_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "byrecovery")
	require.NoError(t, repo.CreateUser(user, ""))

	expiresAt := time.Now().Add(1 * time.Hour)
	const code = "recovery-code-xyz"
	require.NoError(t, repo.SetRecoveryCode(user.ID, code, expiresAt))

	// Pass current time: query finds codes where expires_at >= now
	got, err := repo.GetUserByRecovery(user.Email, code, time.Now())
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestGetUserByRecovery_Expired(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "byrecovery_expired")
	require.NoError(t, repo.CreateUser(user, ""))

	pastExpiry := time.Now().Add(-1 * time.Hour)
	require.NoError(t, repo.SetRecoveryCode(user.ID, "exp-code", pastExpiry))

	// Pass current time: code expired in the past, so expires_at < now → not found
	_, err := repo.GetUserByRecovery(user.Email, "exp-code", time.Now())
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}
