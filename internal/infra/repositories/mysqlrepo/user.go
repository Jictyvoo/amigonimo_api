package mysqlrepo

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
	"github.com/jictyvoo/amigonimo_api/pkg/dberrs"
)

func (r RepoMySQL) CreateUser(user entities.User) (*entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	if user.ID.IsEmpty() {
		user.ID = entities.HexID(uuid.New())
	}

	_, err := r.queries.CreateUser(
		ctx, dbgen.CreateUserParams{
			ID:       user.ID[:],
			Fullname: user.FullName,
			Email:    user.Email,
		},
	)
	if err != nil {
		return nil, WrapError(err, "create user")
	}

	return &user, nil
}

func (r RepoMySQL) GetUserByID(id entities.HexID) (*entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	row, err := r.queries.GetUserByID(ctx, id[:])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dberrs.NewErrDatabaseNotFound("user", id.String(), err)
		}
		return nil, WrapError(err, "get user by id")
	}

	return &entities.User{
		ID:        entities.HexID(row.ID),
		FullName:  row.Fullname,
		UserBasic: entities.UserBasic{Email: row.Email},
	}, nil
}

func (r RepoMySQL) GetUserByEmail(email string) (*entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dberrs.NewErrDatabaseNotFound("user", email, err)
		}
		return nil, WrapError(err, "get user by email")
	}

	return &entities.User{
		ID:        entities.HexID(row.ID),
		FullName:  row.Fullname,
		UserBasic: entities.UserBasic{Email: row.Email},
	}, nil
}
