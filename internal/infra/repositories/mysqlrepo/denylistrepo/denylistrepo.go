package denylistrepo

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
)

var _ denylist.Repository = (*RepoMySQL)(nil)

type RepoMySQL struct {
	mysqlrepo.RepoMySQL
}

func NewRepoMySQL(repoMySQL mysqlrepo.RepoMySQL) *RepoMySQL {
	return &RepoMySQL{RepoMySQL: repoMySQL}
}

func (r *RepoMySQL) AddDenyListEntry(
	p entities.Participant, deniedUserID entities.HexID,
) (entities.DeniedUser, error) {
	id, err := entities.NewHexID()
	if err != nil {
		return entities.DeniedUser{}, err
	}

	ctx, cancel := r.Ctx()
	defer cancel()

	if !p.ID.IsEmpty() {
		_, err = r.Queries().AddDenyListEntryByID(
			ctx, dbgen.AddDenyListEntryByIDParams{
				ID:            id[:],
				ParticipantID: p.ID[:],
				DeniedUserID:  deniedUserID[:],
			},
		)
	} else {
		_, err = r.Queries().AddDenyListEntry(
			ctx, dbgen.AddDenyListEntryParams{
				ID:             id[:],
				UserID:         p.RelatedUser.ID[:],
				SecretFriendID: p.SecretFriendID[:],
				DeniedUserID:   deniedUserID[:],
			},
		)
	}
	if err != nil {
		return entities.DeniedUser{}, mysqlrepo.WrapError(err, "add denylist entry")
	}

	return entities.DeniedUser{
		ID: id,
		Timestamp: entities.Timestamp{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		DeniedUsers: entities.Participant{RelatedUser: entities.User{ID: deniedUserID}},
	}, nil
}

func (r *RepoMySQL) RemoveDenyListEntry(p entities.Participant, deniedUserID entities.HexID) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	var err error
	if !p.ID.IsEmpty() {
		err = r.Queries().RemoveDenyListEntryByID(
			ctx, dbgen.RemoveDenyListEntryByIDParams{
				ParticipantID: p.ID[:],
				DeniedUserID:  deniedUserID[:],
			},
		)
	} else {
		err = r.Queries().RemoveDenyListEntry(
			ctx, dbgen.RemoveDenyListEntryParams{
				UserID:         p.RelatedUser.ID[:],
				SecretFriendID: p.SecretFriendID[:],
				DeniedUserID:   deniedUserID[:],
			},
		)
	}
	if err != nil {
		return mysqlrepo.WrapError(err, "remove denylist entry")
	}

	return nil
}

func (r *RepoMySQL) GetDenyListByParticipant(
	p entities.Participant,
) ([]entities.DeniedUser, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	var err error
	var rows []dbgen.GetDenyListByParticipantRow

	if !p.ID.IsEmpty() {
		var dbRows []dbgen.GetDenyListByParticipantIDRow
		dbRows, err = r.Queries().GetDenyListByParticipantID(ctx, p.ID[:])
		if err == nil {
			rows = make([]dbgen.GetDenyListByParticipantRow, len(dbRows))
			for i, row := range dbRows {
				rows[i] = dbgen.GetDenyListByParticipantRow(row)
			}
		}
	} else {
		rows, err = r.Queries().GetDenyListByParticipant(
			ctx, dbgen.GetDenyListByParticipantParams{
				UserID:         p.RelatedUser.ID[:],
				SecretFriendID: p.SecretFriendID[:],
			},
		)
	}

	if err != nil {
		return nil, mysqlrepo.WrapError(err, "get denylist by participant")
	}

	deniedUsers := make([]entities.DeniedUser, len(rows))
	for i, row := range rows {
		deniedUsers[i] = entities.DeniedUser{
			ID: mappers.HexIDFromBytes(row.ID),
			DeniedUsers: entities.Participant{
				RelatedUser: entities.User{
					ID:       mappers.HexIDFromBytes(row.DeniedUserID),
					FullName: row.Fullname,
					UserBasic: entities.UserBasic{
						Email:    row.Email,
						Username: row.Username,
					},
				},
			},
		}
	}

	return deniedUsers, nil
}
