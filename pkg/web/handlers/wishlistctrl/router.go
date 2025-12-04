package wishlistctrl

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
	handlers := NewWishlistHandlers()

	fuego.Get(server, "/", handlers.GetWishlist)
	fuego.Post(server, "/", handlers.CreateWishlistItem)
	fuego.Delete(server, "/{itemId}", handlers.DeleteWishlistItem)

	return nil
}

func (r *Router) GroupName() string {
	return "/{id}/wishlist"
}

func (r *Router) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
