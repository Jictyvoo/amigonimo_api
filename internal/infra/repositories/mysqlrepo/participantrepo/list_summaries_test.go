//go:build integration

package participantrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/authuserepo"
)

func TestListParticipantSummaries_Success(t *testing.T) {
	repo := genParticipantRepo(t)
	ev := seedUserAndSF(t, repo, "list_sum")
	base := mysqlrepo.NewRepoMySQL(repo.Connection())

	// Add a second user and participant
	secondID, err := entities.NewHexID()
	require.NoError(t, err)
	second := entities.User{
		ID: secondID,
		UserBasic: authvalues.UserBasic{
			Username: "list_sum_second",
			Email:    "list_sum_second@example.com",
			Password: "hashed",
		},
	}
	require.NoError(t, authuserepo.NewRepoMySQL(base).CreateUser(second, ""))

	_, err = repo.AddParticipant(ev.sf.ID, ev.user.ID)
	require.NoError(t, err)
	_, err = repo.AddParticipant(ev.sf.ID, secondID)
	require.NoError(t, err)

	summaries, err := repo.ListParticipantSummaries(ev.sf.ID)
	require.NoError(t, err)
	assert.Len(t, summaries, 2)
}
