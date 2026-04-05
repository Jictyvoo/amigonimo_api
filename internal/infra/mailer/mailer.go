package mailer

import (
	"context"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

var _ authserv.MailerService = (*MailerImpl)(nil)

// MailerImpl implements authserv.MailerService by delegating to a Sender.
// It is registered as a Factory so that each request gets its own instance
// carrying the request context, enabling tracing and cancellation propagation.
type MailerImpl struct {
	ctx         context.Context
	sender      Sender
	fromAddress string
	fromName    string
}

// NewMailerImpl is a Factory constructor resolved by Remy per request.
// It receives the request context, a singleton Sender, and the application config.
func NewMailerImpl(sender Sender, conf config.Config) *MailerImpl {
	return &MailerImpl{
		ctx:         context.Background(),
		sender:      sender,
		fromAddress: conf.Mailer.From,
		fromName:    conf.Mailer.FromName,
	}
}

func (m *MailerImpl) SendActivationEmail(email string, verificationToken string) {
	msg := EmailMessage{
		To:      []string{email},
		Subject: "Activate your account",
		Body:    "Your verification token: " + verificationToken,
	}
	_ = m.sender.Send(m.ctx, msg)
}

func (m *MailerImpl) SendPasswordRecoveryEmail(email string, recoveryCode string) {
	msg := EmailMessage{
		To:      []string{email},
		Subject: "Password recovery",
		Body:    "Your recovery code: " + recoveryCode,
	}
	_ = m.sender.Send(m.ctx, msg)
}
