package authuserepo

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r RepoMySQL) CreateUser(user entities.User, token string) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	if user.ID.IsEmpty() {
		user.ID = entities.HexID(uuid.New())
	}

	_, err := r.Queries().CreateUser(
		ctx, dbgen.CreateUserParams{
			ID:               user.ID[:],
			Fullname:         user.FullName,
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
