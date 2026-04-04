//go:build integration

package wishlistrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
)

func TestRemoveWishlistItem_Success(t *testing.T) {
	repo := genWishlistRepo(t)
	ev := seedMinimalEvent(t, repo, "remove_item")

	ref := wishlist.ParticipantRef{ParticipantID: ev.participant.ID}
	added, err := repo.AddWishlistItem(ref, wishlist.WishlistItem{Label: "To remove"})
	require.NoError(t, err)

	err = repo.RemoveWishlistItem(added.ID, ref)
	require.NoError(t, err)

	items, err := repo.GetWishlistByParticipant(ref)
	require.NoError(t, err)
	assert.Empty(t, items)
}
