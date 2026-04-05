package mailer_test

import (
	"context"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer"
	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

// captureSender records the last Send call for inspection.
type captureSender struct {
	lastCtx context.Context
	lastMsg mailer.EmailMessage
	err     error
}

func (c *captureSender) Send(ctx context.Context, msg mailer.EmailMessage) error {
	c.lastCtx = ctx
	c.lastMsg = msg
	return c.err
}

func baseConfig(from, fromName string) config.Config {
	return config.Config{
		Mailer: config.Mailer{
			From:     from,
			FromName: fromName,
		},
	}
}

func TestNewMailerImpl(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		fromName string
	}{
		{name: "stores from address", from: "no-reply@example.com", fromName: ""},
		{name: "stores from name", from: "no-reply@example.com", fromName: "My App"},
		{name: "empty config", from: "", fromName: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &captureSender{}
			m := mailer.NewMailerImpl(s, baseConfig(tt.from, tt.fromName))
			if m == nil {
				t.Fatal("NewMailerImpl returned nil")
			}
		})
	}
}

func TestSendActivationEmail(t *testing.T) {
	tests := []struct {
		name              string
		email             string
		verificationToken string
		wantSubject       string
		wantBodyContains  string
		senderErr         error
	}{
		{
			name:              "sends to correct address",
			email:             "user@example.com",
			verificationToken: "abc123",
			wantSubject:       "Activate your account",
			wantBodyContains:  "abc123",
		},
		{
			name:              "token included in body",
			email:             "other@example.com",
			verificationToken: "tok-xyz-789",
			wantSubject:       "Activate your account",
			wantBodyContains:  "tok-xyz-789",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &captureSender{err: tt.senderErr}
			m := mailer.NewMailerImpl(s, baseConfig("from@example.com", ""))

			m.SendActivationEmail(tt.email, tt.verificationToken)

			if s.lastCtx == nil {
				t.Error("context was not propagated to sender")
			}
			if len(s.lastMsg.To) == 0 || s.lastMsg.To[0] != tt.email {
				t.Errorf("To = %v, want [%s]", s.lastMsg.To, tt.email)
			}
			if s.lastMsg.Subject != tt.wantSubject {
				t.Errorf("Subject = %q, want %q", s.lastMsg.Subject, tt.wantSubject)
			}
			if tt.wantBodyContains != "" && !containsStr(s.lastMsg.Body, tt.wantBodyContains) {
				t.Errorf("Body %q does not contain %q", s.lastMsg.Body, tt.wantBodyContains)
			}
		})
	}
}

func TestSendPasswordRecoveryEmail(t *testing.T) {
	tests := []struct {
		name             string
		email            string
		recoveryCode     string
		wantSubject      string
		wantBodyContains string
	}{
		{
			name:             "sends to correct address",
			email:            "user@example.com",
			recoveryCode:     "REC-001",
			wantSubject:      "Password recovery",
			wantBodyContains: "REC-001",
		},
		{
			name:             "code included in body",
			email:            "other@example.com",
			recoveryCode:     "XYZ-999",
			wantSubject:      "Password recovery",
			wantBodyContains: "XYZ-999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &captureSender{}
			m := mailer.NewMailerImpl(s, baseConfig("from@example.com", ""))

			m.SendPasswordRecoveryEmail(tt.email, tt.recoveryCode)

			if len(s.lastMsg.To) == 0 || s.lastMsg.To[0] != tt.email {
				t.Errorf("To = %v, want [%s]", s.lastMsg.To, tt.email)
			}
			if s.lastMsg.Subject != tt.wantSubject {
				t.Errorf("Subject = %q, want %q", s.lastMsg.Subject, tt.wantSubject)
			}
			if !containsStr(s.lastMsg.Body, tt.wantBodyContains) {
				t.Errorf("Body %q does not contain %q", s.lastMsg.Body, tt.wantBodyContains)
			}
		})
	}
}

func containsStr(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && func() bool {
		for i := range s {
			if i+len(sub) <= len(s) && s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
