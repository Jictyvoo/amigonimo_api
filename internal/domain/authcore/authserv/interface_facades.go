package authserv

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/checkrecovery"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/login"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/lookuprecovery"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/regenerate"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/requestrecovery"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/resetpassword"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/signup"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/verify"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type MailerService interface {
	signup.Mailer
	requestrecovery.Mailer
}

type UserRetriever interface {
	GetUserByAuthToken(token string) (entities.User, error)
}

type TokenRepository interface {
	UserRetriever
	login.TokenRepository
	regenerate.Repository
}

type UserAuthRepository interface {
	signup.Repository
	login.UserRepository
	verify.Repository
	requestrecovery.UserRepository
	lookuprecovery.Repository
	checkrecovery.Repository
	resetpassword.Repository
}
