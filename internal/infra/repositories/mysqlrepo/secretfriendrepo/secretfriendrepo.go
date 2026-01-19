package secretfriendrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/invite"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

var (
	_ secretfriend.Repository = (*RepoMySQL)(nil)
	_ drawfriends.Repository  = (*RepoMySQL)(nil)
	_ invite.Repository       = (*RepoMySQL)(nil)
)

type RepoMySQL struct {
	mysqlrepo.RepoMySQL
}

func NewRepoMySQL(repoMySQL mysqlrepo.RepoMySQL) *RepoMySQL {
	return &RepoMySQL{RepoMySQL: repoMySQL}
}

func (r *RepoMySQL) CreateSecretFriend(sf *entities.SecretFriend) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	if sf.ID.IsEmpty() {
		var err error
		if sf.ID, err = entities.NewHexID(); err != nil {
			return err
		}
	}

	_, err := r.Queries().CreateSecretFriend(
		ctx, dbgen.CreateSecretFriendParams{
			ID:              sf.ID[:],
			Name:            sf.Name,
			Datetime:        mappers.TimeToNullTime(sf.Datetime),
			Location:        sql.NullString{String: sf.Location, Valid: sf.Location != ""},
			InviteCode:      sf.InviteCode,
			Status:          string(sf.Status),
			OwnerID:         sf.OwnerID[:],
			MaxDenyListSize: sf.MaxDenyListSize,
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "create secret friend")
	}

	return nil
}

func (r *RepoMySQL) GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbSF, err := r.Queries().GetSecretFriendByID(ctx, id[:])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.SecretFriend{}, dberrs.NewErrDatabaseNotFound(
				"secret_friend", id.String(), err,
			)
		}
		return entities.SecretFriend{}, mysqlrepo.WrapError(err, "get secret friend by id")
	}

	sf := mappers.ToEntitySecretFriend(dbSF)

	// Fetch participants
	// TODO: Optimize to only fetch partipants if required
	participants, err := r.ListParticipants(ctx, id)
	if err == nil {
		sf.Participants = participants
	}

	return sf, nil
}

func (r *RepoMySQL) UpdateSecretFriend(sf *entities.SecretFriend) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	err := r.Queries().UpdateSecretFriend(
		ctx, dbgen.UpdateSecretFriendParams{
			ID:       sf.ID[:],
			Name:     sf.Name,
			Datetime: mappers.TimeToNullTime(sf.Datetime),
			Location: sql.NullString{String: sf.Location, Valid: sf.Location != ""},
			Status:   string(sf.Status),
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "update secret friend")
	}

	return nil
}

func (r *RepoMySQL) ListParticipants(
	ctx context.Context,
	sfID entities.HexID,
) ([]entities.Participant, error) {
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

func (r *RepoMySQL) GetDrawResultForUser(
	secretFriendID, userID entities.HexID,
) (entities.DrawResultItem, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	row, err := r.Queries().GetDrawResultForUser(
		ctx, dbgen.GetDrawResultForUserParams{
			SecretFriendID: secretFriendID[:],
			UserID:         userID[:],
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.DrawResultItem{}, dberrs.NewErrDatabaseNotFound(
				"draw_result",
				userID.String(),
				err,
			)
		}
		return entities.DrawResultItem{}, mysqlrepo.WrapError(err, "get draw result for user")
	}

	return entities.DrawResultItem{
		Timestamp: entities.Timestamp{
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		},
		Giver: entities.Participant{
			ID: mappers.HexIDFromBytes(row.GiverParticipantID),
			RelatedUser: entities.User{
				ID: mappers.HexIDFromBytes(row.GiverUserID),
			},
		},
		Receiver: entities.Participant{
			ID: mappers.HexIDFromBytes(row.ReceiverParticipantID),
			RelatedUser: entities.User{
				ID:       mappers.HexIDFromBytes(row.ReceiverUserID),
				FullName: row.ReceiverFullname,
				UserBasic: entities.UserBasic{
					Email: row.ReceiverEmail,
				},
			},
		},
	}, nil
}

func (r *RepoMySQL) SaveDrawResults(
	secretFriendID entities.HexID,
	results []entities.DrawResultItem,
) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	// Use a transaction since we are saving multiple results
	txFinisher, err := r.BeginTx(ctx, nil)
	if err != nil {
		return mysqlrepo.WrapError(err, "begin transaction for save draw results")
	}
	defer func() {
		_ = txFinisher(err == nil)
	}()

	for _, result := range results {
		_, err = r.Queries().SaveDrawResult(
			ctx, dbgen.SaveDrawResultParams{
				GiverParticipantID:    result.Giver.ID[:],
				ReceiverParticipantID: result.Receiver.ID[:],
				SecretFriendID:        secretFriendID[:],
			},
		)
		if err != nil {
			return mysqlrepo.WrapError(
				err,
				fmt.Sprintf("save draw result for giver %s", result.Giver.ID.String()),
			)
		}
	}

	return nil
}

func (r *RepoMySQL) ListSecretFriends(
	userID entities.HexID,
) ([]entities.SecretFriend, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	rows, err := r.Queries().
		ListSecretFriends(ctx, dbgen.ListSecretFriendsParams{UserID: userID[:], OwnerID: userID[:]})
	if err != nil {
		return nil, mysqlrepo.WrapError(err, "list secret friends")
	}

	sfList := make([]entities.SecretFriend, len(rows))
	for i, row := range rows {
		sfList[i] = mappers.ToEntitySecretFriend(row)
	}

	return sfList, nil
}

func (r *RepoMySQL) GetSecretFriendByInviteCode(code string) (entities.SecretFriend, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbSF, err := r.Queries().GetSecretFriendByInviteCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.SecretFriend{}, dberrs.NewErrDatabaseNotFound(
				"secret_friend",
				code,
				err,
			)
		}
		return entities.SecretFriend{}, mysqlrepo.WrapError(err, "get secret friend by invite code")
	}

	return mappers.ToEntitySecretFriend(dbSF), nil
}
