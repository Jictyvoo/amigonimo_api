package main

import (
	"errors"
	"net/http"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
)

func main() {
	conf := bootstrap.Config()
	db := bootstrap.OpenDatabase(conf.Database)
	defer func() {
		_ = db.Close()
	}()

	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	remy.RegisterInstance(inj, db)
	bootstrap.DoInjections(inj, conf)

	// Parse RSA key and extract public key for JWT middleware
	jwtPublicKey, secretErr := registerSecret([]byte(conf.Runtime.AuthSecretKey), inj)
	if secretErr != nil {
		panic(secretErr)
	}
	conf.Runtime.AuthSecretKey = "" // Empty the secret key after injection

	// Create web server
	//goland:noinspection GoResourceLeak
	server, err := bootstrap.NewWebServer(conf, jwtPublicKey, inj)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = server.Close()
	}()
	if err = server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
