package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const envConfFile = "CONF_FILE"

// Load loads the configuration using filename as the TOML file path.
// If the CONF_FILE environment variable is set it takes precedence over
// the filename argument. A missing file is not treated as an error.
// Environment variables are always overlaid on top of the loaded file.
func Load(filename string) (Config, error) {
	confPath := os.Getenv(envConfFile)
	if confPath == "" {
		confPath = filename
	}

	conf, loadErr := LoadTOML(confPath)
	if loadErr != nil && !os.IsNotExist(errors.Unwrap(loadErr)) {
		return Config{}, fmt.Errorf("load %s: %w", confPath, loadErr)
	}

	if err := LoadConfigFromEnv(&conf); err != nil {
		return Config{}, fmt.Errorf("load env: %w", err)
	}
	return conf, nil
}

func LoadConfigFromEnv(config *Config) error {
	var useDebugStr string

	err := BindEnv(
		BindField(&useDebugStr, envUseDebug, strings.ToLower),
		BindField(&config.ProjectName, envProjectName, strings.TrimSpace),
		BindField(&config.Runtime.APILocale, envAPILocale, strings.TrimSpace),
		BindField(&config.Runtime.Host, envAPIHost, strings.ToLower),
		BindField(&config.Runtime.Port, envAPIPort, nil),
		BindField(&config.Runtime.AuthSecretKey, envAuthSecretKey, strings.TrimSpace),
	)
	if err != nil {
		return err
	}

	// Handle IsDebug special logic
	if !config.IsDebug && useDebugStr != "" {
		config.IsDebug = useDebugStr != "false"
	}

	return LoadDatabaseFromEnv(&config.Database)
}

func LoadDatabaseFromEnv(conf *Database) error {
	return BindEnv(
		BindField(&conf.Host, envDatabaseHost, strings.ToLower),
		BindField(&conf.Port, envDatabasePort, nil),
		BindField(&conf.User, envDatabaseUser, nil),
		BindField(&conf.Password, envDatabasePassword, nil),
		BindField(&conf.Database, envDatabaseName, nil),
		BindFieldErr(&conf.Timeout, envDatabaseTimeout, time.ParseDuration),
	)
}
