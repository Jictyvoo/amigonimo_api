package config

const (
	envUseDebug    = "DEBUG"
	envProjectName = "PROJECT_NAME"

	envDatabaseHost     = "DATABASE_HOST"
	envDatabasePort     = "DATABASE_PORT"
	envDatabaseUser     = "DATABASE_USER"
	envDatabasePassword = "DATABASE_PASSWORD"
	envDatabaseName     = "DATABASE_NAME"

	envAPILocale     = "API_LOCALE"
	envAPIHost       = "API_HOST"
	envAPIPort       = "API_PORT"
	envAuthSecretKey = "AUTH_SECRET_KEY"
)

func DefaultConfig() Config {
	return Config{
		IsDebug:     false,
		ProjectName: "amigonymus_api",
		Runtime: Runtime{
			APILocale: "pt_BR",
			Port:      8649,
		},
		Database: Database{
			Host:     "localhost",
			Port:     3306,
			User:     "secretshhh",
			Password: "testing_u-know",
			Database: "amigonimo_db",
		},
	}
}
