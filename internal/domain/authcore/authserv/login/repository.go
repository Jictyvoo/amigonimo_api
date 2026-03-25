package login

import "github.com/jictyvoo/amigonimo_api/internal/entities"

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=user_repository_mock_test.go -package=login github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/login UserRepository
//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=token_repository_mock_test.go -package=login github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/login TokenRepository

type UserRepository interface {
	GetUserByEmailOrUsername(email, username string) (entities.User, error)
}

type TokenRepository interface {
	GetAuthenticationToken(userID entities.HexID) (entities.AuthenticationToken, error)
	UpsertAuthToken(authentication *entities.AuthenticationToken) error
}
