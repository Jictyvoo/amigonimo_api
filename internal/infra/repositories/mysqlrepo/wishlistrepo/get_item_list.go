package wishlistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func (r *RepoMySQL) GetWishlistByParticipant(
	participant entities.Participant,
) ([]entities.WishlistItem, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(participant)
	rows, err := r.Queries().GetWishlistByParticipant(
		ctx, dbgen.GetWishlistByParticipantParams{
			ParticipantID:  ids.participantID,
			UserID:         ids.userID,
			SecretFriendID: ids.secretFriendID,
		},
	)
	if err != nil {
		return nil, mysqlrepo.WrapError(err, "get wishlist by participant")
	}

	items := make([]entities.WishlistItem, len(rows))
	for i, row := range rows {
		items[i] = mappers.ToEntityWishlistItem(row)
	}

	return items, nil
}
