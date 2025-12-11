package main

import (
	"errors"
	"net/http"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/dashboardctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/invitesctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
)

func main() {
	conf := bootstrap.Config()
	db := bootstrap.OpenDatabase(conf.Database)
	defer db.Close()

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
	server, err := web.NewServer(
		conf, jwtPublicKey, web.WithPublicRouters(
			authctrl.NewAuthRouter(
				authctrl.Config{
					ActiveRoutes: authctrl.RouteLogin | authctrl.RouteSignUp | authctrl.RouteRegenerateToken | authctrl.RouteResetPassword,
					Injector:     inj,
				},
			),
		),
		web.WithPrivateRouters(
			dashboardctrl.NewRouter(),
			invitesctrl.NewRouter(),
			secretfriendsctrl.NewRouter(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err = server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
