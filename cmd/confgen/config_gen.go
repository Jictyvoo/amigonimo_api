package main

import (
	"bytes"
	"errors"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
)

func configGenCMD() (err error) {
	// check if config file exists and if not create it
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
	}

	if os.IsNotExist(err) {
		file, err = os.OpenFile(config.DefaultFileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	} else {
		// update config file to have all needed values
		file, err = os.OpenFile(config.DefaultFileName, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	}

	// check if error occurred and then write default config file
	if err == nil {
		if configData.Runtime.AuthSecretKey == "" {
			configData.Runtime.AuthSecretKey = generateAuthKey()
		}
		marshaledData, _ := toml.Marshal(&configData)
		marshaledData = bytes.TrimSpace(marshaledData)
		_, err = file.Write(marshaledData)
		_, _ = file.WriteString("\n")
	}

	if err != nil {
		return
	}
	err = file.Close()
	return
}
