package denylistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func (r *RepoMySQL) GetDenyListByParticipant(
	participant denylist.ParticipantRef,
) ([]denylist.DeniedEntry, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(participant)
	rows, err := r.Queries().GetDenyListByParticipant(
		ctx, dbgen.GetDenyListByParticipantParams{
			ParticipantID:  ids.participantID,
			UserID:         ids.userID,
			SecretFriendID: ids.secretFriendID,
		},
	)
	if err != nil {
		return nil, mysqlrepo.WrapError(err, "get denylist by participant")
	}

	entries := make([]denylist.DeniedEntry, len(rows))
	for i, row := range rows {
		entries[i] = mappers.ToDeniedEntry(row)
	}

	return entries, nil
}
