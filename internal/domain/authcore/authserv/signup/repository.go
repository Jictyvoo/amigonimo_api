package signup

import "github.com/jictyvoo/amigonimo_api/internal/entities"

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=signup github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/signup Repository
//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=mailer_mock_test.go -package=signup github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/signup Mailer

type Repository interface {
	GetUserByEmailOrUsername(email, username string) (entities.User, error)
	CreateUser(user entities.User, token string) error
}

type Mailer interface {
	SendActivationEmail(email string, verificationToken string)
}
