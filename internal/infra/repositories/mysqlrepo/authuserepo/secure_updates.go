package authuserepo

import (
	"database/sql"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r RepoMySQL) SetUserVerified(userID entities.HexID) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	err := r.Queries().SetUserVerified(ctx, userID[:])
	if err != nil {
		return mysqlrepo.WrapError(err, "set user verified")
	}

	return nil
}

func (r RepoMySQL) SetRecoveryCode(userID entities.HexID, code string, expiresAt time.Time) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	err := r.Queries().SetRecoveryCode(
		ctx, dbgen.SetRecoveryCodeParams{
			RecoveryCode:          sql.NullString{String: code, Valid: code != ""},
			RecoveryCodeExpiresAt: sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()},
			ID:                    userID[:],
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "set recovery code")
	}

	return nil
}

func (r RepoMySQL) UpdatePassword(userID entities.HexID, newPassword string) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	err := r.Queries().UpdatePassword(
		ctx, dbgen.UpdatePasswordParams{
			Password: newPassword,
			ID:       userID[:],
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "update password")
	}

	return nil
}
