package wishlistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r *RepoMySQL) RemoveWishlistItem(
	itemID entities.HexID,
	participant entities.Participant,
) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(participant)
	err := r.Queries().RemoveWishlistItem(
		ctx, dbgen.RemoveWishlistItemParams{
			ID:             itemID[:],
			ParticipantID:  ids.participantID,
			UserID:         ids.userID,
			SecretFriendID: ids.secretFriendID,
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "remove wishlist item")
	}

	return nil
}
