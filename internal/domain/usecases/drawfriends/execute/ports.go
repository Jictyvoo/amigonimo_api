package execute

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=ports_mock_test.go -package=execute github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/execute Repository,secretFriendFacadePort

type Repository interface {
	dbrock.Transactioner

	SaveDrawResults(secretFriendID entities.HexID, results []entities.DrawResultItem) error
}

type secretFriendFacadePort interface {
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateStatus(id entities.HexID, status entities.SecretFriendStatus) error
}

type SecretFriendFacade interface {
	ports.Facade
	secretFriendFacadePort
}
