package verify

import "github.com/jictyvoo/amigonimo_api/internal/entities"

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=verify github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/verify Repository

type Repository interface {
	GetUserByVerificationCode(code string) (entities.User, error)
	SetUserVerified(userID entities.HexID) error
}
