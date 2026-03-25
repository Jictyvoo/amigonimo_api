package changeusername

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/session"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=user_repository_mock_test.go -package=changeusername github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeusername UserRepository
//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=changeusername github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeusername Repository

type UserRepository interface {
	GetUserByUsername(username string) (entities.User, error)
}

type Repository interface {
	session.Repository
	UpdateUsername(userID entities.HexID, username string) error
}
