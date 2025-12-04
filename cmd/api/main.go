package main

import (
	"errors"
	"net/http"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/bootstrap"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/dashboardctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/drawresultctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/invitesctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
)

func main() {
	conf := bootstrap.Config()
	db := bootstrap.OpenDatabase(conf.Database)
	defer db.Close()

	inj := remy.NewInjector(remy.Config{DuckTypeElements: true})
	remy.RegisterInstance(inj, db)
	bootstrap.DoInjections(inj, conf)

	// Create web server
	server := web.NewServer(
		conf, web.WithPublicRouters(authctrl.NewRouter()),
		web.WithPrivateRouters(
			dashboardctrl.NewRouter(),
			denylistctrl.NewRouter(),
			drawresultctrl.NewRouter(),
			invitesctrl.NewRouter(),
			participantsctrl.NewRouter(),
			secretfriendsctrl.NewRouter(),
			wishlistctrl.NewRouter(),
		),
	)

	if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
