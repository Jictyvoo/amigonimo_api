package authtokenrepo

import (
	"database/sql"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/internal/dbgen"
)

func (r RepoMySQL) UpsertAuthToken(authentication *entities.AuthenticationToken) error {
	ctx, cancel := r.Ctx()
	defer cancel()

	// If ID is empty, generate a new ID for the token
	if authentication.ID.IsEmpty() {
		newID, err := entities.NewHexID()
		if err != nil {
			return err
		}
		authentication.ID = newID
	}

	var refreshToken sql.NullString
	if authentication.RefreshToken.Valid {
		refreshToken = sql.NullString{
			String: authentication.RefreshToken.UUID.String(),
			Valid:  true,
		}
	}

	_, err := r.Queries().UpsertAuthToken(
		ctx, dbgen.UpsertAuthTokenParams{
			ID:           authentication.ID[:],
			UserID:       authentication.User.ID[:],
			Token:        authentication.AuthToken,
			RefreshToken: refreshToken,
			ExpiresAt: sql.NullTime{
				Time:  authentication.ExpiresAt,
				Valid: !authentication.ExpiresAt.IsZero(),
			},
		},
	)
	if err != nil {
		return mysqlrepo.WrapError(err, "upsert authentication token")
	}

	return nil
}
