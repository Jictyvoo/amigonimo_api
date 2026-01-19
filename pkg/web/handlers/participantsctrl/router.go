package participantsctrl

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant"
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
	var useCaseFac UseCaseFactory[participant.UseCase] = func(ctx context.Context) (participant.UseCase, error) {
		return remy.GetWithContext[participant.UseCase](r.injector, ctx)
	}
	handlers := NewParticipantsHandlers(useCaseFac)

	fuego.Post(server, "/", handlers.ConfirmParticipation, web.OptionAuthToken())
	fuego.Get(server, "/", handlers.ListParticipants)

	return nil
}

func (r *Router) GroupName() string {
	return "/participants"
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
