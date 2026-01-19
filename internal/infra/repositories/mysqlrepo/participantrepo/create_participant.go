package participantrepo

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r *RepoMySQL) AddParticipant(sfID, userID entities.HexID) (entities.Participant, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	id, err := entities.NewHexID()
	if err != nil {
		return entities.Participant{}, err
	}
	_, err = r.Queries().AddParticipant(
		ctx, dbgen.AddParticipantParams{
			ID:             id[:],
			SecretFriendID: sfID[:],
			UserID:         userID[:],
		},
	)
	if err != nil {
		return entities.Participant{}, mysqlrepo.WrapError(err, "add participant")
	}

	now := time.Now()
	return entities.Participant{
		Timestamp: entities.Timestamp{
			CreatedAt: now,
			UpdatedAt: now,
		},
		ID:             id,
		RelatedUser:    entities.User{ID: userID}, // TODO: Retrieve username/fullname
		SecretFriendID: sfID,
		JoinedAt:       now,
	}, nil
}
