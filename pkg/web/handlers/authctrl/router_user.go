package authctrl

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"

	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
)

type RouterUser struct {
	activeRoutes DefinedRoute
	middlewares  []web.HttpMiddleware
}

func NewUserRouter(activeRoutes DefinedRoute) *RouterUser {
	return &RouterUser{middlewares: []web.HttpMiddleware{}, activeRoutes: activeRoutes}
}

func (r *RouterUser) GroupName() string {
	return "users/edit"
}

func (r *RouterUser) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}

func (r *RouterUser) RegisterRoutes(server *fuego.Server) error {
	ctrl := controllers.NewUserEditionController(nil, nil, nil)
	// Bind handlers
	groupTag := option.Tags("User Management")
	if r.activeRoutes.is(RouteEditPassword) {
		fuego.Patch(
			server, "/password", ctrl.EditUserPassword,
			option.Summary("Edit user's password"),
			option.Description(
				"Edit a user's password by providing an authentication token and a new password",
			),
			groupTag, web.OptionAuthToken(),
		)
	}

	if r.activeRoutes.is(RouteEditUsername) {
		fuego.Patch(
			server, "/username", ctrl.EditUsername,
			option.Summary("Edit username"),
			option.Description(
				"Edit a user's username by providing an authentication token and a new username",
			),
			groupTag, web.OptionAuthToken(),
		)
	}

	if r.activeRoutes.is(RouteEditEmail) {
		fuego.Patch(
			server, "/email", ctrl.EditUserEmail,
			option.Summary("Edit User email and send a new verification code in the email"),
			option.Description(
				"Edit a user's email address and send a new verification code by providing an authentication token and a new email",
			),
			groupTag, web.OptionAuthToken(),
		)
	}

	return nil
}

func (r *RouterUser) AddMiddleware(middleware web.HttpMiddleware) {
	r.middlewares = append(r.middlewares, middleware)
}
