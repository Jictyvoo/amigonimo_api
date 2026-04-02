package sqltest

import (
	"errors"
	"fmt"
	"os"
	"unicode"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

// loadConfig loads conf.toml from the current working directory if present,
// then overlays any DATABASE_* / DEBUG / etc. environment variables on top.
// A missing conf.toml is not an error; env vars alone are sufficient.
func loadConfig() (config.Config, error) {
	conf, loadErr := config.LoadTOML(config.DefaultFileName)
	if loadErr != nil && !os.IsNotExist(errors.Unwrap(loadErr)) {
		return config.Config{}, fmt.Errorf("load %s: %w", config.DefaultFileName, loadErr)
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
