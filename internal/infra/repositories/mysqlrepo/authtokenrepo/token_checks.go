package authtokenrepo

import (
	"database/sql"
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/mappers"
	"github.com/jictyvoo/amigonimo_api/pkg/dberrs"
)

func (r RepoMySQL) GetUserByAuthToken(token string) (entities.User, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbUser, err := r.Queries().GetUserByAuthToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, dberrs.NewErrDatabaseNotFound("user", token, err)
		}
		return entities.User{}, mysqlrepo.WrapError(err, "get user by auth token")
	}

	return mappers.ToEntityUser(dbUser), nil
}

func (r RepoMySQL) GetAuthenticationToken(
	userID entities.HexID,
) (entities.AuthenticationToken, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	dbToken, err := r.Queries().GetAuthenticationToken(ctx, userID[:])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.AuthenticationToken{
				ID:   entities.HexID{},
				User: entities.User{ID: userID},
			}, nil
		}
		return entities.AuthenticationToken{}, mysqlrepo.WrapError(err, "get authentication token")
	}

	// Note: User info should be populated by the service layer
	// For now, return token with minimal user info (service has the full user)
	return mappers.ToEntityAuthenticationToken(dbToken), nil
}

func (r RepoMySQL) CheckAuthenticationByRefreshToken(
	authToken string,
) (*entities.AuthenticationToken, error) {
	ctx, cancel := r.Ctx()
	defer cancel()

	refreshToken := sql.NullString{
		String: authToken,
		Valid:  true,
	}

	dbToken, err := r.Queries().CheckAuthenticationByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dberrs.NewErrDatabaseNotFound("authentication token", authToken, err)
		}
		return nil, mysqlrepo.WrapError(err, "check authentication by refresh token")
	}

	// Note: User info should be populated by the service layer
	// For now, return token with minimal user info
	token := mappers.ToEntityAuthenticationToken(dbToken)
	return &token, nil
}
