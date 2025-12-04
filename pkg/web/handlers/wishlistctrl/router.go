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

func (r *Router) V1(server *fuego.Server) error {
	handlers := NewWishlistHandlers()

	fuego.Get(server, "/{id}/wishlist", handlers.GetWishlist)
	fuego.Post(server, "/{id}/wishlist", handlers.CreateWishlistItem)
	fuego.Delete(server, "/{id}/wishlist/{itemId}", handlers.DeleteWishlistItem)

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
