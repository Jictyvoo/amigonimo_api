package denylistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r *RepoMySQL) RemoveDenyListEntry(p entities.Participant, deniedUserID entities.HexID) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	ids := r.queryIDs(p)
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
