package sqltest

import (
	"errors"
	"fmt"
	"os"
	"unicode"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

const envConfFile = "CONF_FILE"

// loadConfig loads the configuration TOML file and overlays env vars on top.
// File resolution order:
//  1. Path given by the CONF_FILE environment variable (explicit override)
//  2. conf.toml in the current working directory
//  3. Default values only (env vars are always applied last)
func loadConfig() (config.Config, error) {
	confPath := os.Getenv(envConfFile)
	if confPath == "" {
		confPath = config.DefaultFileName
	}

	conf, loadErr := config.LoadTOML(confPath)
	if loadErr != nil && !os.IsNotExist(errors.Unwrap(loadErr)) {
		return config.Config{}, fmt.Errorf("load %s: %w", confPath, loadErr)
	}

	// Allow env vars to override individual fields
	if err := config.LoadConfigFromEnv(&conf); err != nil {
		return config.Config{}, fmt.Errorf("load env: %w", err)
	}
	return conf, nil
}

// dbNameNormalizer converts a test name into a valid MySQL identifier, prefixed
// with "test_". Non-alphanumeric characters are replaced with underscores.
func dbNameNormalizer(name string) string {
	runes := []rune(name)
	for i, ch := range runes {
		ch = unicode.ToLower(ch)
		if !unicode.IsLetter(ch) && !unicode.IsNumber(ch) {
			ch = '_'
		}
		runes[i] = ch
	}
	return "sqlrepo_" + string(runes)
}
