package checkrecovery

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=checkrecovery github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/checkrecovery Repository

type Repository interface {
	GetUserByRecovery(userEmail string, code string, expiredAt time.Time) (entities.User, error)
}
