package participantrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
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
