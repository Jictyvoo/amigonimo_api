package requestrecovery

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=user_repository_mock_test.go -package=requestrecovery github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/requestrecovery UserRepository
//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=mailer_mock_test.go -package=requestrecovery github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/requestrecovery Mailer

type UserRepository interface {
	GetUserByEmail(userEmail string) (entities.User, error)
	SetRecoveryCode(userID entities.HexID, code string, expiresAt time.Time) error
}

type Mailer interface {
	SendPasswordRecoveryEmail(email string, recoveryCode string)
}
