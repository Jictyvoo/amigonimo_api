package config

import "time"

const DefaultFileName = "conf.toml"

type (
	Runtime struct {
		APILocale     string `toml:"api_locale"`
		Host          string `toml:"host"`
		Port          uint16 `toml:"port"`
		AuthSecretKey string `toml:"auth_secret_key"`
	}

	Database struct {
		Host     string        `toml:"host"`
		Port     uint16        `toml:"port"`
		User     string        `toml:"user"`
		Password string        `toml:"password"`
		Database string        `toml:"database"`
		Timeout  time.Duration `toml:"timeout"`
	}

	MailerSmtp struct {
		Host       string `toml:"host"`
		Port       uint16 `toml:"port"`
		User       string `toml:"user"`
		Password   string `toml:"password"`
		Encryption string `toml:"encryption"` // "tls", "ssl", or "none"
	}

	MailerWebhook struct {
		URL    string `toml:"url"`
		APIKey string `toml:"api_key"`
	}

	Mailer struct {
		Driver   string        `toml:"driver"` // "smtp", "webhook", or "stub"
		From     string        `toml:"from"`
		FromName string        `toml:"from_name"`
		Smtp     MailerSmtp    `toml:"smtp"`
		Webhook  MailerWebhook `toml:"webhook"`
	}

	Config struct {
		IsDebug     bool     `toml:"is_debug"`
		ProjectName string   `toml:"project_name"`
		Runtime     Runtime  `toml:"server"`
		Database    Database `toml:"database"`
		Mailer      Mailer   `toml:"mailer"`
	}
)
