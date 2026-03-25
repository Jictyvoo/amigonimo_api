package resetpassword

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=resetpassword github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/resetpassword Repository

type Repository interface {
	GetUserByRecovery(userEmail string, code string, expiredAt time.Time) (entities.User, error)
	UpdatePassword(userID entities.HexID, newPassword string) error
	SetRecoveryCode(userID entities.HexID, code string, expiresAt time.Time) error
}
