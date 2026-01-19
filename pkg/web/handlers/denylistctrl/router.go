package denylistctrl

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type Router struct {
	middlewares []web.HttpMiddleware
	injector    remy.Injector
}

func NewRouter(injector remy.Injector) *Router {
	return &Router{middlewares: []web.HttpMiddleware{}, injector: injector}
}

func (r *Router) RegisterRoutes(server *fuego.Server) error {
	var useCaseFac UseCaseFactory[denylist.UseCase] = func(ctx context.Context) (denylist.UseCase, error) {
		return remy.GetWithContext[denylist.UseCase](r.injector, ctx)
	}
	handlers := NewController(useCaseFac)

	fuego.Get(server, "/", handlers.GetDenyList, web.OptionAuthToken())
	fuego.Post(server, "/", handlers.AddDenyListEntry, web.OptionAuthToken())
	fuego.Delete(server, "/{deniedUserId}", handlers.RemoveDenyListEntry, web.OptionAuthToken())

	return nil
}

func (r *Router) GroupName() string {
	return "/denylist"
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
