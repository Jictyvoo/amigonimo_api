package web

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
)

type HttpMiddleware = func(http.Handler) http.Handler

type (
	RouterWithSubRouter interface {
		SubRouters() []RouterContract
	}
	RouterMiddlewareExtender interface {
		AddMiddleware(middleware HttpMiddleware)
	}
	RouterContract interface {
		RegisterRoutes(server *fuego.Server) error
		GroupName() string
		Middlewares() []HttpMiddleware
	}
)

func SetupRoutes(server *fuego.Server, routerList ...RouterContract) error {
	for _, router := range routerList {
		groupName := router.GroupName()
		groupRouter := fuego.Group(
			server, groupName,
			fuego.OptionMiddleware(router.Middlewares()...),
		)

		// Setup versioned route
		if err := router.RegisterRoutes(groupRouter); err != nil {
			return fmt.Errorf("failed to setup routes: %w", err)
		}

		if withSubRouter, ok := router.(RouterWithSubRouter); ok {
			if err := SetupRoutes(groupRouter, withSubRouter.SubRouters()...); err != nil {
				return err
			}
		}
	}

	return nil
}
