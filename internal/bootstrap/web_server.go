package bootstrap

import (
	"crypto/rsa"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/pkg/config"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
)

func NewWebServer(
	conf config.Config, jwtPublicKey *rsa.PublicKey, inj remy.Injector,
) (*web.Server, error) {
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
			secretfriendsctrl.NewRouter(inj),
		),
	)
	return server, err
}
