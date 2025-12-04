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

	fuego.Get(server, "/{id}/denylist", handlers.GetDenyList)
	fuego.Post(server, "/{id}/denylist", handlers.AddDenyListEntry)
	fuego.Delete(server, "/{id}/denylist/{targetUserId}", handlers.RemoveDenyListEntry)

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
