package bootstrap

import (
	"log/slog"
	"os"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func Config() config.Config {
	conf, err := config.Load(config.DefaultFileName)
	if err != nil {
		slog.Error("Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return conf
}
