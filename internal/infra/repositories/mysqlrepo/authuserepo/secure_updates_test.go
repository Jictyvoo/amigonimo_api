//go:build integration

package authuserepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUserVerified_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "setverified")
	require.NoError(t, repo.CreateUser(user, "token"))

	err := repo.SetUserVerified(user.ID)
	assert.NoError(t, err)
}

func TestSetRecoveryCode_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "setrecovery")
	require.NoError(t, repo.CreateUser(user, ""))

	expiresAt := time.Now().Add(1 * time.Hour)
	err := repo.SetRecoveryCode(user.ID, "recovery-abc", expiresAt)
	require.NoError(t, err)

	// Pass current time: code is still valid (expires in future)
	got, err := repo.GetUserByRecovery(user.Email, "recovery-abc", time.Now())
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestUpdatePassword_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "updatepass")
	require.NoError(t, repo.CreateUser(user, ""))

	err := repo.UpdatePassword(user.ID, "new_hashed_password")
	require.NoError(t, err)

	got, err := repo.GetUserByEmail(user.Email)
	require.NoError(t, err)
	assert.Equal(t, "new_hashed_password", got.Password)
}
