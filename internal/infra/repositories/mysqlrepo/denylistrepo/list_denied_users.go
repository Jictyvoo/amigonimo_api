package denylistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func (r *RepoMySQL) GetDenyListByParticipant(
	p entities.Participant,
) ([]entities.DeniedUser, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(p)
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

	deniedUsers := make([]entities.DeniedUser, len(rows))
	for i, row := range rows {
		deniedUsers[i] = mappers.ToEntityDeniedUser(row)
	}

	return deniedUsers, nil
}
