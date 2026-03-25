package lookuprecovery

import "github.com/jictyvoo/amigonimo_api/internal/entities"

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=lookuprecovery github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/lookuprecovery Repository

type Repository interface {
	GetUserByUsername(username string) (entities.User, error)
}
