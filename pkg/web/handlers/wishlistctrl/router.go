package wishlistctrl

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/wishlist"
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
	var useCaseFac UseCaseFactory[wishlist.UseCase] = func(ctx context.Context) (wishlist.UseCase, error) {
		return remy.GetWithContext[wishlist.UseCase](r.injector, ctx)
	}
	handlers := NewController(useCaseFac)

	fuego.Get(server, "/", handlers.GetWishlist, web.OptionAuthToken())
	fuego.Post(server, "/", handlers.CreateWishlistItem, web.OptionAuthToken())
	fuego.Delete(server, "/{itemId}", handlers.DeleteWishlistItem, web.OptionAuthToken())

	return nil
}

func (r *Router) GroupName() string {
	return "/wishlist"
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
