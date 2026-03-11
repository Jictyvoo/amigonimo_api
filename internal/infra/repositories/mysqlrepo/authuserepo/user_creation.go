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

	_, err := r.Queries().CreateUser(
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

	return nil
}
