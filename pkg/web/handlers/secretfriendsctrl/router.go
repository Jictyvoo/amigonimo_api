package secretfriendsctrl

import (
	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
)

type Router struct {
	middlewares []web.HttpMiddleware
}

func NewRouter() *Router {
	return &Router{middlewares: []web.HttpMiddleware{}}
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

func (r *Router) SubRouters() []web.RouterContract {
	return []web.RouterContract{
		wishlistctrl.NewRouter(),
		denylistctrl.NewRouter(),
		participantsctrl.NewRouter(),
	}
}
