package mailfacade

import "github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"

var _ authserv.MailerService = (*MailerImpl)(nil)

type MailerImpl struct{}

func NewMailerImpl() *MailerImpl {
	return &MailerImpl{}
}

func (m MailerImpl) SendActivationEmail(email string, verificationToken string) {
	// TODO implement me
	panic("implement me")
}

func (m MailerImpl) SendPasswordRecoveryEmail(email string, recoveryCode string) {
	// TODO implement me
	panic("implement me")
}
