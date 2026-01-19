package secretfriendsctrl

import (
	"github.com/go-fuego/fuego"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
)

type Router struct {
	middlewares []web.HttpMiddleware
	injector    remy.Injector
}

func NewRouter(inj remy.Injector) *Router {
	return &Router{injector: inj, middlewares: []web.HttpMiddleware{}}
}

func (r *Router) RegisterRoutes(server *fuego.Server) error {
	handlers := NewSecretFriendsHandlers()

	fuego.Post(server, "", handlers.CreateSecretFriend)
	fuego.Get(server, "/{id}", handlers.GetSecretFriend)
	fuego.Patch(server, "/{id}", handlers.UpdateSecretFriend)
	fuego.Post(server, "/{id}/draw", handlers.DrawSecretFriend)
	fuego.Get(server, "/{id}/draw-result", handlers.GetDrawResult)

	return nil
}

func (r *Router) GroupName() string {
	return "/secret-friends"
}

func (r *Router) AddMiddleware(middleware web.HttpMiddleware) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}

func (r *Router) SubRouters() (subgroupPattern string, routers []web.RouterContract) {
	return "/{id}", []web.RouterContract{
		wishlistctrl.NewRouter(r.injector),
		denylistctrl.NewRouter(r.injector),
		participantsctrl.NewRouter(r.injector),
	}
}
