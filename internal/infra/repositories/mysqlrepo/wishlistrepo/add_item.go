package wishlistrepo

import (
	"database/sql"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r *RepoMySQL) AddWishlistItem(
	participant entities.Participant,
	wishItem entities.WishlistItem,
) (entities.WishlistItem, error) {
	id, err := entities.NewHexID()
	if err != nil {
		return entities.WishlistItem{}, err
	}

	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(participant)
	_, err = r.Queries().AddWishlistItem(
		ctx, dbgen.AddWishlistItemParams{
			ID:             id[:],
			ParticipantID:  ids.participantID,
			UserID:         ids.userID,
			SecretFriendID: ids.secretFriendID,
			Label:          wishItem.Label,
			Comments: sql.NullString{
				String: wishItem.Comments,
				Valid:  wishItem.Comments != "",
			},
		},
	)
	if err != nil {
		return entities.WishlistItem{}, mysqlrepo.WrapError(err, "add wishlist item")
	}

	wishItem.ID = id
	return wishItem, nil
}
