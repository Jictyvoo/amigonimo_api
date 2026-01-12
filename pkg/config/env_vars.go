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
)

const (
	defaultAPIPort      = 8649
	defaultDatabasePort = 3306
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
	}
}
