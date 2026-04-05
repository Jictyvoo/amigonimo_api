package smtp

import (
	"context"
	netsmtp "net/smtp"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
)

var _ mailer.Sender = (*Sender)(nil)

// Sender delivers email via STARTTLS SMTP.
type Sender struct {
	config Config
}

func New(cfg Config) *Sender {
	return &Sender{config: cfg}
}

func (s *Sender) Send(ctx context.Context, msg mailer.EmailMessage) error {
	parts := &emailParts{
		sender:   s.config.SenderHeader(),
		to:       msg.To,
		subject:  msg.Subject,
		body:     msg.Body,
		bodyHTML: msg.BodyHTML,
	}

	for _, a := range msg.Attachments {
		att := attachment{
			contentType: a.ContentType,
			filename:    a.Filename,
			data:        a.Data,
		}
		if a.Inline {
			parts.htmlAttachments = append(parts.htmlAttachments, att)
		} else {
			parts.attachments = append(parts.attachments, att)
		}
	}

	raw, err := buildMail(ctx, parts)
	if err != nil {
		return err
	}

	auth := netsmtp.PlainAuth("", s.config.User, s.config.Password, s.config.Host)
	return netsmtp.SendMail(s.config.Address(), auth, s.config.FromAddress, msg.To, raw)
}
