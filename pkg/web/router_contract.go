package web

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
)

type HttpMiddleware = func(http.Handler) http.Handler

type (
	RouterWithSubRouter interface {
		SubRouters() (subgroupPattern string, routers []RouterContract)
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

		if withSubRouter, ok := router.(RouterWithSubRouter); ok {
			subGroupPattern, subRouterList := withSubRouter.SubRouters()
			subGroupRouter := groupRouter
			if subGroupPattern != "" {
				subGroupRouter = fuego.Group(groupRouter, subGroupPattern)
			}
			if err := SetupRoutes(subGroupRouter, subRouterList...); err != nil {
				return err
			}
		}

		// Setup versioned route
		if err := router.RegisterRoutes(groupRouter); err != nil {
			return fmt.Errorf("failed to register `%s` routes: %w", groupName, err)
		}
	}

	return nil
}
