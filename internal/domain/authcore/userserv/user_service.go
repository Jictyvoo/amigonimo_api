package userserv

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeemail"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changepassword"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeusername"
)

type UserEditionService struct {
	ChangePassword changepassword.UseCase
	ChangeEmail    changeemail.UseCase
	ChangeUsername changeusername.UseCase
}

func NewUserEditService(
	userRepository authserv.UserAuthRepository,
	userEditRepository UserEditionRepository,
	mailer authserv.MailerService,
) UserEditionService {
	return UserEditionService{
		ChangePassword: changepassword.New(userRepository, userEditRepository),
		ChangeEmail:    changeemail.New(userRepository, userEditRepository, mailer),
		ChangeUsername: changeusername.New(userRepository, userEditRepository),
	}
}
