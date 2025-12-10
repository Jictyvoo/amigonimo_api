package authuserepo

import "github.com/jictyvoo/amigonimo_api/internal/entities"

func (r RepoMySQL) SetUserVerified(userID entities.HexID) error {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) SetRecoveryCode(userID entities.HexID, code string) error {
	// TODO implement me
	panic("implement me")
}

func (r RepoMySQL) UpdatePassword(userID entities.HexID, newPassword string) error {
	// TODO implement me
	panic("implement me")
}
