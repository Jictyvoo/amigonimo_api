package mailer

import "context"

// Driver identifies the email transport backend.
type Driver string

const (
	DriverSMTP    Driver = "smtp"
	DriverWebhook Driver = "webhook"
	DriverStub    Driver = "stub"
)

// Attachment is a file attached to an outgoing email.
// Set Inline to true for HTML-embedded resources (images, etc.).
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
	Inline      bool
}

// EmailMessage holds the data needed to send an email,
// independent of the transport backend.
type EmailMessage struct {
	To          []string
	Subject     string
	Body        string
	BodyHTML    string
	Attachments []Attachment
}

// Sender is the composition seam between MailerImpl and the transport backend.
// ctx is propagated to allow tracing and cancellation at the transport level.
type Sender interface {
	Send(ctx context.Context, msg EmailMessage) error
}
