package authuserepo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func (r RepoMySQL) GetUserByUsername(username string) (entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbUser, err := r.Queries().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, dberrs.NewErrDatabaseNotFound("user", username, err)
		}
		return entities.User{}, mysqlrepo.WrapError(err, "get user by username")
	}

	return mappers.ToEntityUser(dbUser), nil
}

func (r RepoMySQL) GetUserByVerificationCode(code string) (entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	verificationCode := sql.NullString{
		String: code,
		Valid:  true,
	}

	dbUser, err := r.Queries().GetUserByVerificationCode(ctx, verificationCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, dberrs.NewErrDatabaseNotFound("user", code, err)
		}
		return entities.User{}, mysqlrepo.WrapError(err, "get user by verification code")
	}

	return mappers.ToEntityUser(dbUser), nil
}

func (r RepoMySQL) GetUserByRecovery(
	userEmail string, code string, expiredAt time.Time,
) (entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbUser, err := r.Queries().GetUserByRecovery(
		ctx, dbgen.GetUserByRecoveryParams{
			Email:                 userEmail,
			RecoveryCode:          sql.NullString{String: code, Valid: true},
			RecoveryCodeExpiresAt: sql.NullTime{Time: expiredAt, Valid: true},
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, dberrs.NewErrDatabaseNotFound("user", userEmail, err)
		}
		return entities.User{}, mysqlrepo.WrapError(err, "get user by recovery")
	}

	return mappers.ToEntityUser(dbUser), nil
}

func (r RepoMySQL) GetUserByEmailOrUsername(email, username string) (entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbUser, err := r.Queries().GetUserByEmailOrUsername(
		ctx, dbgen.GetUserByEmailOrUsernameParams{
			Email:    email,
			Username: username,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, dberrs.NewErrDatabaseNotFound(
				"user", email+" or "+username, err,
			)
		}
		return entities.User{}, mysqlrepo.WrapError(err, "get user by email or username")
	}

	return mappers.ToEntityUser(dbUser), nil
}

func (r RepoMySQL) GetUserByEmail(email string) (decodedUser entities.User, err error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	row, err := r.Queries().GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return decodedUser, dberrs.NewErrDatabaseNotFound("user", email, err)
		}
		return decodedUser, mysqlrepo.WrapError(err, "get user by email")
	}

	return mappers.ToEntityUser(row), nil
}
