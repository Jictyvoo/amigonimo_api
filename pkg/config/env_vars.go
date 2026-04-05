package config

import "time"

const (
	envUseDebug    = "DEBUG"
	envProjectName = "PROJECT_NAME"

	envDatabaseHost     = "DATABASE_HOST"
	envDatabasePort     = "DATABASE_PORT"
	envDatabaseUser     = "DATABASE_USER"
	envDatabasePassword = "DATABASE_PASSWORD"
	envDatabaseName     = "DATABASE_NAME"
	envDatabaseTimeout  = "DATABASE_TIMEOUT"

	envAPILocale     = "API_LOCALE"
	envAPIHost       = "API_HOST"
	envAPIPort       = "API_PORT"
	envAuthSecretKey = "AUTH_SECRET_KEY" //nolint:gosec

	envMailerDriver   = "MAILER_DRIVER"
	envMailerFrom     = "MAILER_FROM"
	envMailerFromName = "MAILER_FROM_NAME"

	envSmtpHost       = "SMTP_HOST"
	envSmtpPort       = "SMTP_PORT"
	envSmtpUser       = "SMTP_USER"
	envSmtpPassword   = "SMTP_PASSWORD" //nolint:gosec
	envSmtpEncryption = "SMTP_ENCRYPTION"

	envWebhookURL    = "MAILER_WEBHOOK_URL"
	envWebhookAPIKey = "MAILER_WEBHOOK_API_KEY" //nolint:gosec
)

const (
	defaultAPIPort      = 8649
	defaultDatabasePort = 3306
	defaultSmtpPort     = 587
)

func DefaultConfig() Config {
	return Config{
		IsDebug:     false,
		ProjectName: "amigonymus_api",
		Runtime: Runtime{
			APILocale: "pt_BR",
			Port:      defaultAPIPort,
		},
		Database: Database{
			Host:     "localhost",
			Port:     defaultDatabasePort,
			User:     "secretshhh",
			Password: "testing_u-know",
			Database: "amigonimo_db",
			Timeout:  time.Second,
		},
		Mailer: Mailer{
			Driver: "stub",
			Smtp: MailerSmtp{
				Port:       defaultSmtpPort,
				Encryption: "tls",
			},
		},
	}
}
