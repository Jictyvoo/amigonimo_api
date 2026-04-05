package stub_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer/stub"
)

func TestStubSend(t *testing.T) {
	tests := []struct {
		name string
		msg  mailer.EmailMessage
	}{
		{
			name: "plain text email",
			msg: mailer.EmailMessage{
				To:      []string{"user@example.com"},
				Subject: "Hello",
				Body:    "World",
			},
		},
		{
			name: "html email",
			msg: mailer.EmailMessage{
				To:       []string{"user@example.com"},
				Subject:  "Hello HTML",
				BodyHTML: "<p>World</p>",
			},
		},
		{
			name: "email with attachments",
			msg: mailer.EmailMessage{
				To:      []string{"user@example.com"},
				Subject: "With attachment",
				Body:    "See attached",
				Attachments: []mailer.Attachment{
					{Filename: "file.txt", ContentType: "text/plain", Data: []byte("data")},
				},
			},
		},
		{
			name: "multiple recipients",
			msg: mailer.EmailMessage{
				To:      []string{"a@example.com", "b@example.com"},
				Subject: "Broadcast",
				Body:    "Hi all",
			},
		},
		{
			name: "empty message",
			msg:  mailer.EmailMessage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stub.New()
			err := s.Send(context.Background(), tt.msg)
			if err != nil {
				t.Errorf("Send() returned unexpected error: %v", err)
			}
		})
	}
}

func TestStubSendContextPropagation(t *testing.T) {
	type ctxKey struct{}
	ctx := context.WithValue(context.Background(), ctxKey{}, "value")

	s := stub.New()
	// Stub always returns nil regardless of context, but should not panic.
	err := s.Send(ctx, mailer.EmailMessage{To: []string{"u@example.com"}, Subject: "test"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStubImplementsSender(t *testing.T) {
	var _ mailer.Sender = stub.New()
}

func TestStubNeverReturnsError(t *testing.T) {
	s := stub.New()
	// Even with a cancelled context the stub must not error.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.Send(ctx, mailer.EmailMessage{})
	if errors.Is(err, context.Canceled) {
		t.Error("stub should not return context errors")
	}
}
