//go:build integration

package authuserepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

func makeTestUser(t testing.TB, suffix string) entities.User {
	t.Helper()
	id, err := entities.NewHexID()
	require.NoError(t, err)
	return entities.User{
		ID: id,
		UserBasic: authvalues.UserBasic{
			Username: "testuser_" + suffix,
			Email:    "testuser_" + suffix + "@example.com",
			Password: "hashed_password",
		},
	}
}

func TestCreateUser_Success(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "create_ok")

	err := repo.CreateUser(user, "verification-token")
	require.NoError(t, err)

	got, err := repo.GetUserByEmail(user.Email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Email, got.Email)
	assert.Equal(t, user.Username, got.Username)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo := genAuthUserRepo(t)
	user := makeTestUser(t, "dup_email")

	require.NoError(t, repo.CreateUser(user, "token1"))
	err := repo.CreateUser(user, "token2")
	assert.Error(t, err)
}
