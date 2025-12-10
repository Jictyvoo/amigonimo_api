package authtokenrepo

import "github.com/jictyvoo/amigonimo_api/internal/entities"

func (r RepoMySQL) GetUserByAuthToken(token string) (entities.User, error) {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) GetAuthenticationToken(
	userID entities.HexID,
) (entities.AuthenticationToken, error) {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) CheckAuthenticationByRefreshToken(
	authToken string,
) (*entities.AuthenticationToken, error) {
	// TODO implement me
	panic("implement me")
}
