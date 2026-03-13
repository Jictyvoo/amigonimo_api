package denylistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r *RepoMySQL) RemoveDenyListEntry(
	participant denylist.ParticipantRef,
	deniedUserID entities.HexID,
) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(participant)
	err := r.Queries().RemoveDenyListEntry(
		ctx, dbgen.RemoveDenyListEntryParams{
			ParticipantID:  ids.participantID,
			UserID:         ids.userID,
			SecretFriendID: ids.secretFriendID,
			DeniedUserID:   deniedUserID[:],
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "remove denylist entry")
	}

	return nil
}
