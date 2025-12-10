package authuserepo

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/pkg/dberrs"
)

func (r RepoMySQL) GetUserByUsername(username string) (entities.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) GetUserByVerificationCode(code string) (entities.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) GetUserByRecovery(
	userEmail string,
	code string,
	expiredAt time.Time,
) (entities.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) GetUserByEmailOrUsername(email, username string) (entities.User, error) {
	// TODO implement me
	panic("implement me")
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

	return entities.User{
		ID:        entities.HexID(row.ID),
		FullName:  row.Fullname,
		UserBasic: entities.UserBasic{Email: row.Email},
	}, nil
}
