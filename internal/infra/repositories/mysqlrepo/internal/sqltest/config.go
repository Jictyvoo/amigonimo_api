package sqltest

import (
	"fmt"
	"unicode"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

// loadConfig finds conf.toml by walking up from the current working directory
// to the module root (identified by go.mod), then loads the config.
func loadConfig() (config.Config, error) {
	conf, loadErr := config.LoadTOML(config.DefaultFileName)
	if loadErr != nil {
		return config.Config{}, fmt.Errorf("load %s: %w", config.DefaultFileName, loadErr)
	}

	// Allow env vars to override individual fields
	_ = config.LoadConfigFromEnv(&conf)
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
