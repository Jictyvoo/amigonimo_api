//go:build integration

package secretfriendrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/drawdto"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/secretfriendrepo"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestSaveAndGetDrawResult_Success(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "draw_ok")

	sf := &entities.SecretFriend{
		Name:            "draw-sf",
		InviteCode:      "draw-code",
		Status:          entities.StatusDraft,
		OwnerID:         owner.ID,
		MaxDenyListSize: 3,
	}
	require.NoError(t, repo.CreateSecretFriend(sf))

	p1 := seedParticipant(t, repo, sf.ID, owner.ID)

	// Create a second user and participant
	secondUser := seedOwner(t, repo, "draw_ok_second")
	p2 := seedParticipant(t, repo, sf.ID, secondUser.ID)

	results := []drawdto.DrawResultItem{
		{GiverParticipantID: p1.ID, ReceiverParticipantID: p2.ID},
		{GiverParticipantID: p2.ID, ReceiverParticipantID: p1.ID},
	}

	require.NoError(t, repo.SaveDrawResults(sf.ID, results))

	// SaveDrawResults uses BeginTx which mutates the queries field; use a fresh repo
	// instance sharing the same underlying *sql.DB to avoid stale transaction state.
	freshRepo := secretfriendrepo.NewRepoMySQL(mysqlrepo.NewRepoMySQL(repo.Connection()))
	got, err := freshRepo.GetDrawResultForUser(sf.ID, owner.ID)
	require.NoError(t, err)
	assert.Equal(t, p1.ID, got.GiverParticipantID)
	assert.Equal(t, p2.ID, got.ReceiverParticipantID)
}

func TestGetDrawResultForUser_NotFound(t *testing.T) {
	repo := genSecretFriendRepo(t)
	owner := seedOwner(t, repo, "draw_nf")

	sf := &entities.SecretFriend{
		Name:            "draw-nf-sf",
		InviteCode:      "draw-nf-code",
		Status:          entities.StatusDraft,
		OwnerID:         owner.ID,
		MaxDenyListSize: 3,
	}
	require.NoError(t, repo.CreateSecretFriend(sf))

	_, err := repo.GetDrawResultForUser(sf.ID, owner.ID)
	var notFound *dberrs.ErrDatabaseNotFound
	assert.ErrorAs(t, err, &notFound)
}
