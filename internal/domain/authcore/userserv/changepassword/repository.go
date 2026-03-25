package changepassword

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=changepassword github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changepassword Repository

type Repository interface {
	UpdatePassword(userID entities.HexID, newPassword string) error
}
