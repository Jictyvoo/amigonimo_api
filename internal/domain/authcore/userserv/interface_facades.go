package userserv

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeemail"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changepassword"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/changeusername"
)

type UserEditionRepository interface {
	changepassword.Repository
	changeemail.Repository
	changeusername.Repository
}
