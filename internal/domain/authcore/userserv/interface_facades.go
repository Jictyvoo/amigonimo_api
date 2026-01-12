package userserv

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=../../mocks/useredit_repo_mock.go -package=mocks github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv UserEditionRepository

type UserEditionRepository interface {
	GetUserByAuthCode(authToken string) (entities.User, error)
	ChangeEmail(userID entities.HexID, newEmail string) error
	SetNewVerificationCode(userID entities.HexID, code string) error
	UpdateUsername(userID entities.HexID, username string) error
}
