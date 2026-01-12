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

	Config struct {
		IsDebug     bool     `toml:"is_debug"`
		ProjectName string   `toml:"project_name"`
		Runtime     Runtime  `toml:"server"`
		Database    Database `toml:"database"`
	}
)
