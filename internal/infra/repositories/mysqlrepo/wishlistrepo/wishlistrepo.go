package wishlistrepo

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
)

var _ wishlist.Repository = (*RepoMySQL)(nil)

type RepoMySQL struct {
	mysqlrepo.RepoMySQL
}

func NewRepoMySQL(repoMySQL mysqlrepo.RepoMySQL) *RepoMySQL {
	return &RepoMySQL{RepoMySQL: repoMySQL}
}

func (r *RepoMySQL) queryIDs(participant entities.Participant) (
	ids struct {
		participantID  []byte
		userID         []byte
		secretFriendID []byte
	},
) {
	if !participant.ID.IsEmpty() {
		ids.participantID = participant.ID[:]
		return ids
	}

	ids.userID = participant.RelatedUser.ID[:]
	ids.secretFriendID = participant.SecretFriendID[:]

	return ids
}
