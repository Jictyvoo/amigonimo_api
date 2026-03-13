package denylistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
)

var _ denylist.Repository = (*RepoMySQL)(nil)

type RepoMySQL struct {
	mysqlrepo.RepoMySQL
}

func NewRepoMySQL(repoMySQL mysqlrepo.RepoMySQL) *RepoMySQL {
	return &RepoMySQL{RepoMySQL: repoMySQL}
}

func (r *RepoMySQL) queryIDs(participant denylist.ParticipantRef) (
	ids struct {
		participantID  []byte
		userID         []byte
		secretFriendID []byte
	},
) {
	if !participant.ParticipantID.IsEmpty() {
		ids.participantID = participant.ParticipantID[:]
		return ids
	}

	ids.userID = participant.UserID[:]
	ids.secretFriendID = participant.SecretFriendID[:]

	return ids
}
