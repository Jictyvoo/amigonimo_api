package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

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

func LoadTOML(paths ...string) (config Config, err error) {
	config = DefaultConfig()

	for _, path := range paths {
		var (
			file      *os.File
			fileBytes []byte
		)

		cleanPath := filepath.Clean(path)
		file, err = os.OpenFile(cleanPath, os.O_RDONLY, 0)
		if os.IsNotExist(err) || file == nil {
			err = errors.Wrap(err, fmt.Sprintf("Unable to load config from %s", path))
			return
		}

		if fileBytes, err = io.ReadAll(file); err == nil {
			err = toml.Unmarshal(fileBytes, &config)
			if err != nil {
				return
			}
		}
	}

	return
}
