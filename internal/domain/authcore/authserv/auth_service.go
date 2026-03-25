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
)

type AuthService struct {
	SignUp                  signup.UseCase
	LogIn                   login.UseCase
	RegenerateToken         regenerate.UseCase
	VerifyUser              verify.UseCase
	RequestPasswordRecovery requestrecovery.UseCase
	LookupRecoveryContact   lookuprecovery.UseCase
	CheckRecoveryCode       checkrecovery.UseCase
	ResetPassword           resetpassword.UseCase
}

func NewAuthService(
	userRepository UserAuthRepository,
	tokenRepository TokenRepository,
	mailer MailerService,
) AuthService {
	return AuthService{
		SignUp:                  signup.New(userRepository, mailer),
		LogIn:                   login.New(userRepository, tokenRepository),
		RegenerateToken:         regenerate.New(tokenRepository),
		VerifyUser:              verify.New(userRepository),
		RequestPasswordRecovery: requestrecovery.New(userRepository, mailer),
		LookupRecoveryContact:   lookuprecovery.New(userRepository),
		CheckRecoveryCode:       checkrecovery.New(userRepository),
		ResetPassword:           resetpassword.New(userRepository),
	}
}
