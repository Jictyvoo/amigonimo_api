package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func configGenCMD() (err error) {
	// check if a config file exists and if not create it
	var (
		file       *os.File
		configData config.Config
	)

	configData, err = config.LoadTOML(config.DefaultFileName)
	if err != nil {
		err = errors.Unwrap(err)
		if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
			err = unwrappedErr
		}

		if !os.IsNotExist(err) {
			return err
		}
	}

	const filePerm = 0o600
	if file, err = os.OpenFile(config.DefaultFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, filePerm); err != nil {
		return fmt.Errorf("open config file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(file)

	// check if an error occurred and then write the default config to the file
	if configData.Runtime.AuthSecretKey == "" {
		configData.Runtime.AuthSecretKey = generateAuthKey()
	}

	if err = toml.NewEncoder(file).Encode(configData); err != nil {
		return fmt.Errorf("failed to encode config using toml: %w", err)
	}
	return nil
}
