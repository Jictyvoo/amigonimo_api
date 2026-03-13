package denylistrepo

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r *RepoMySQL) AddDenyListEntry(
	participant denylist.ParticipantRef, deniedUserID entities.HexID,
) (entities.DeniedUser, error) {
	id, err := entities.NewHexID()
	if err != nil {
		return entities.DeniedUser{}, err
	}

	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(participant)
	_, err = r.Queries().AddDenyListEntry(
		ctx, dbgen.AddDenyListEntryParams{
			ID:             id[:],
			ParticipantID:  ids.participantID,
			UserID:         ids.userID,
			SecretFriendID: ids.secretFriendID,
			DeniedUserID:   deniedUserID[:],
		},
	)
	if err != nil {
		return entities.DeniedUser{}, mysqlrepo.WrapError(err, "add denylist entry")
	}

	return entities.DeniedUser{
		ID: id,
		Timestamp: entities.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		InnerParticipant: entities.Participant{RelatedUser: entities.User{ID: deniedUserID}},
	}, nil
}
