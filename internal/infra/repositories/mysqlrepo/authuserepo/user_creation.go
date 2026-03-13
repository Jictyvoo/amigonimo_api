package authuserepo

import (
	"database/sql"
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r RepoMySQL) CreateUser(user entities.User, token string) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	if user.ID.IsEmpty() {
		return errors.New("user ID is required")
	}

	profileID, err := entities.NewHexID()
	if err != nil {
		return err
	}

	onFinishTx, err := r.BeginTx(ctx, nil)
	if err != nil {
		return mysqlrepo.WrapError(err, "begin tx for create user")
	}

	committed := false
	defer func() {
		_ = onFinishTx(committed)
	}()

	_, err = r.Queries().CreateUser(
		ctx, dbgen.CreateUserParams{
			ID:               user.ID[:],
			Email:            user.Email,
			Username:         user.Username,
			Password:         user.Password,
			VerificationCode: sql.NullString{String: token, Valid: token != ""},
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "create user")
	}

	if _, err = r.Queries().CreateUserProfile(
		ctx,
		dbgen.CreateUserProfileParams{
			ID:        profileID[:],
			UserID:    user.ID[:],
			Fullname:  sql.NullString{String: user.FullName, Valid: user.FullName != ""},
			Nickname:  sql.NullString{},
			ImageLink: sql.NullString{},
			Birthday:  sql.NullTime{},
			Address:   sql.NullString{},
		},
	); err != nil {
		return mysqlrepo.WrapError(err, "create user profile")
	}

	committed = true
	return nil
}
