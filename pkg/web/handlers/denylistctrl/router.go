package denylistctrl

import (
	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type Router struct {
	middlewares []web.HttpMiddleware
}

func NewRouter() *Router {
	return &Router{middlewares: []web.HttpMiddleware{}}
}

func (r *Router) RegisterRoutes(server *fuego.Server) error {
	handlers := NewDenyListHandlers()

	fuego.Get(server, "/", handlers.GetDenyList)
	fuego.Post(server, "/", handlers.AddDenyListEntry)
	fuego.Delete(server, "/{targetUserId}", handlers.RemoveDenyListEntry)

	return nil
}

func (r *Router) GroupName() string {
	return "/{id}/denylist"
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
