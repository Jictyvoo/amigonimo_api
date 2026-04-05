package stub

import (
	"context"
	"log/slog"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
)

var _ mailer.Sender = (*Sender)(nil)

// Sender is a no-op mailer that logs outgoing messages.
// Intended for development and test environments.
type Sender struct{}

func New() *Sender { return &Sender{} }

func (s *Sender) Send(_ context.Context, msg mailer.EmailMessage) error {
	slog.Info(
		"[stub mailer] sending email",
		slog.Any("to", msg.To),
		slog.String("subject", msg.Subject),
		slog.String("body", msg.Body),
		slog.Int("attachments", len(msg.Attachments)),
	)
	return nil
}
