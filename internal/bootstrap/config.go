package bootstrap

import (
	"errors"
	"log/slog"
	"os"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func Config() config.Config {
	conf, confErr := config.LoadTOML(config.DefaultFileName)
	if confErr != nil {
		if os.IsNotExist(errors.Unwrap(confErr)) {
			slog.Warn(
				"Failed to load config file",
				slog.String("file", config.DefaultFileName),
				slog.String("error", confErr.Error()),
			)
		}
	}

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		slog.Error(
			"Error during environment variables load",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	return conf
}
