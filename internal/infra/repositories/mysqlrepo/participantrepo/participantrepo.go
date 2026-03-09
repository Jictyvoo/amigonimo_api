package participantrepo

import (
	"database/sql"
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

var _ participant.Repository = (*RepoMySQL)(nil)

type RepoMySQL struct {
	mysqlrepo.RepoMySQL
}

func NewRepoMySQL(repoMySQL mysqlrepo.RepoMySQL) *RepoMySQL {
	return &RepoMySQL{RepoMySQL: repoMySQL}
}

func parseAndReturnParticipant(
	err error, userID entities.HexID, dbP dbgen.Participant,
) (entities.Participant, error) {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Participant{}, dberrs.NewErrDatabaseNotFound(
				"participant", userID.String(), err,
			)
		}
		return entities.Participant{}, mysqlrepo.WrapError(err, "get participant by sf and user")
	}

	return mappers.ToEntityParticipant(dbP), nil
}

func (r *RepoMySQL) RemoveParticipant(sfID, userID entities.HexID) error {
	ctx, cancel := r.Ctx()
	defer cancel()
	err := r.Queries().DeleteParticipantBySFAndUser(
		ctx, dbgen.DeleteParticipantBySFAndUserParams{
			SecretFriendID: sfID[:],
			UserID:         userID[:],
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "remove participant")
	}
	return nil
}
