package changeemail

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/session"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=user_repository_mock_test.go -package=changeemail github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeemail UserRepository
//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=changeemail github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeemail Repository
//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=mailer_mock_test.go -package=changeemail github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeemail Mailer

type UserRepository interface {
	GetUserByEmail(email string) (entities.User, error)
}

type Repository interface {
	session.Repository
	ChangeEmail(userID entities.HexID, newEmail string) error
	SetNewVerificationCode(userID entities.HexID, code string) error
}

type Mailer interface {
	SendActivationEmail(email string, verificationToken string)
}
