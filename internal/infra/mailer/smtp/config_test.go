package smtp_test

import (
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/infra/mailer/smtp"
)

func TestConfigAddress(t *testing.T) {
	tests := []struct {
		name string
		cfg  smtp.Config
		want string
	}{
		{
			name: "standard SMTP port",
			cfg:  smtp.Config{Host: "smtp.example.com", Port: 587},
			want: "smtp.example.com:587",
		},
		{
			name: "SSL port",
			cfg:  smtp.Config{Host: "mail.example.com", Port: 465},
			want: "mail.example.com:465",
		},
		{
			name: "localhost",
			cfg:  smtp.Config{Host: "localhost", Port: 25},
			want: "localhost:25",
		},
		{
			name: "zero port",
			cfg:  smtp.Config{Host: "host.example.com", Port: 0},
			want: "host.example.com:0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.Address()
			if got != tt.want {
				t.Errorf("Address() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfigSenderHeader(t *testing.T) {
	tests := []struct {
		name string
		cfg  smtp.Config
		want string
	}{
		{
			name: "with display name",
			cfg:  smtp.Config{FromAddress: "no-reply@example.com", FromName: "My App"},
			want: "My App <no-reply@example.com>",
		},
		{
			name: "without display name",
			cfg:  smtp.Config{FromAddress: "no-reply@example.com", FromName: ""},
			want: "no-reply@example.com",
		},
		{
			name: "name with spaces",
			cfg:  smtp.Config{FromAddress: "bot@acme.io", FromName: "Acme Notifications"},
			want: "Acme Notifications <bot@acme.io>",
		},
		{
			name: "empty config",
			cfg:  smtp.Config{},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.SenderHeader()
			if got != tt.want {
				t.Errorf("SenderHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}
