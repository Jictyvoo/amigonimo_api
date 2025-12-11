package mailfacade

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
)

var _ authserv.MailerService = (*MailerImpl)(nil)

type MailerImpl struct{}

func NewMailerImpl() *MailerImpl {
	return &MailerImpl{}
}

func (m MailerImpl) SendActivationEmail(email string, verificationToken string) {
	fmt.Println(email, verificationToken)
}

func (m MailerImpl) SendPasswordRecoveryEmail(email string, recoveryCode string) {
	fmt.Println(email, recoveryCode)
}
