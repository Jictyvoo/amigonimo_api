package participantrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func (r *RepoMySQL) ListParticipantSummaries(sfID entities.HexID) ([]participant.Summary, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	rows, err := r.Queries().ListParticipantsBySecretFriend(ctx, sfID[:])
	if err != nil {
		return nil, mysqlrepo.WrapError(err, "list participant summaries")
	}

	summaries := make([]participant.Summary, len(rows))
	for i, row := range rows {
		summaries[i] = participant.Summary{
			Participant: mappers.MapParticipantRow(row),
			FullName:    row.Fullname,
		}
	}

	return summaries, nil
}
