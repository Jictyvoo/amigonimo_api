//go:build integration

package wishlistrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
)

func TestGetWishlistByParticipant_Empty(t *testing.T) {
	repo := genWishlistRepo(t)
	ev := seedMinimalEvent(t, repo, "list_empty")

	ref := wishlist.ParticipantRef{ParticipantID: ev.participant.ID}
	items, err := repo.GetWishlistByParticipant(ref)
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestGetWishlistByParticipant_WithItems(t *testing.T) {
	repo := genWishlistRepo(t)
	ev := seedMinimalEvent(t, repo, "list_items")

	ref := wishlist.ParticipantRef{ParticipantID: ev.participant.ID}
	_, err := repo.AddWishlistItem(ref, wishlist.WishlistItem{Label: "Item 1"})
	require.NoError(t, err)
	_, err = repo.AddWishlistItem(ref, wishlist.WishlistItem{Label: "Item 2"})
	require.NoError(t, err)

	items, err := repo.GetWishlistByParticipant(ref)
	require.NoError(t, err)
	assert.Len(t, items, 2)
}
