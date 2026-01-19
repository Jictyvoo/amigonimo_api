package participantrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

func (r *RepoMySQL) GetParticipant(
	sfID, userID entities.HexID,
) (entities.Participant, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbP, err := r.Queries().GetParticipantBySFAndUser(
		ctx, dbgen.GetParticipantBySFAndUserParams{
			SecretFriendID: sfID[:],
			UserID:         userID[:],
		},
	)
	return parseAndReturnParticipant(err, userID, dbP)
}

func (r *RepoMySQL) GetParticipantByID(userID entities.HexID) (entities.Participant, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbP, err := r.Queries().GetParticipantByID(ctx, userID[:])
	return parseAndReturnParticipant(err, userID, dbP)
}

func (r *RepoMySQL) ListParticipants(sfID entities.HexID) ([]entities.Participant, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	rows, err := r.Queries().ListParticipantsBySecretFriend(ctx, sfID[:])
	if err != nil {
		return nil, mysqlrepo.WrapError(err, "list participants")
	}

	participants := make([]entities.Participant, len(rows))
	for i, row := range rows {
		participants[i] = mappers.MapParticipantRow(row)
	}

	return participants, nil
}
