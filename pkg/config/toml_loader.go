package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

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
			err = fmt.Errorf("unable to load config from %s: %w", path, err)
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
