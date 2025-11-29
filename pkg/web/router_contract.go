package web

import (
	"fmt"
	"net/http"

	"github.com/go-fuego/fuego"
)

type HttpMiddleware = func(http.Handler) http.Handler

type RouterContract interface {
	V1(server *fuego.Server) error
	GroupName() string
	Middlewares() []HttpMiddleware
	AddMiddleware(middleware HttpMiddleware)
}

func SetupRoutes(server *fuego.Server, routerList ...RouterContract) error {
	for _, router := range routerList {
		groupName := router.GroupName()
		groupRouter := fuego.Group(
			server, groupName,
			fuego.OptionMiddleware(router.Middlewares()...),
		)

		// Setup versioned route
		if err := router.V1(groupRouter); err != nil {
			return fmt.Errorf("failed to setup routes: %w", err)
		}
	}

	return nil
}
