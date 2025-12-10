package participantsctrl

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
	handlers := NewParticipantsHandlers()

	fuego.Post(server, "/", handlers.ConfirmParticipation)
	fuego.Post(server, "/confirm", handlers.ConfirmParticipation)
	fuego.Get(server, "/", handlers.ListParticipants)

	return nil
}

func (r *Router) GroupName() string {
	return "/{id}/participants"
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
