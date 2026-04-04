//go:build integration

package wishlistrepo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
)

func TestAddWishlistItem_Success(t *testing.T) {
	repo := genWishlistRepo(t)
	ev := seedMinimalEvent(t, repo, "add_item")

	ref := wishlist.ParticipantRef{ParticipantID: ev.participant.ID}
	item := wishlist.WishlistItem{Label: "Book", Comments: "Any edition"}

	got, err := repo.AddWishlistItem(ref, item)
	require.NoError(t, err)
	assert.False(t, got.ID.IsEmpty())
	assert.Equal(t, "Book", got.Label)
	assert.Equal(t, "Any edition", got.Comments)
}
