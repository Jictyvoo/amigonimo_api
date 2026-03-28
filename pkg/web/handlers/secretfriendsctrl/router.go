package secretfriendsctrl

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
)

type Router struct {
	middlewares []web.HttpMiddleware
	injector    remy.Injector
}

func NewRouter(injector remy.Injector) *Router {
	return &Router{middlewares: []web.HttpMiddleware{}, injector: injector}
}

func (r *Router) RegisterRoutes(server *fuego.Server) error {
	var sfFac UseCaseFactory[secretfriend.UseCase] = func(ctx context.Context) (secretfriend.UseCase, error) {
		return remy.GetWithContext[secretfriend.UseCase](r.injector, ctx)
	}
	var drawFac UseCaseFactory[drawfriends.Service] = func(ctx context.Context) (drawfriends.Service, error) {
		return remy.GetWithContext[drawfriends.Service](r.injector, ctx)
	}

	ctrl := NewController(sfFac, drawFac)

	groupTag := option.Tags("Secret Friends")
	optionEventID := option.Path("id", "Secret Friend ID", param.Required())

	fuego.Post(
		server, "/", ctrl.CreateSecretFriend,
		option.Summary("Create a new secret friend event"),
		option.Description(
			"Create a new secret friend event with name, datetime, location, and optional max deny list size",
		),
		groupTag, web.OptionAuthToken(),
	)

	fuego.Get(
		server, "/", ctrl.GetSecretFriendList,
		option.Summary("List all user's linked secret-friends"),
		option.Description(
			"List secret-friends grouped by active/inactive",
		),
		groupTag, web.OptionAuthToken(),
	)

	fuego.Get(
		server, "/invites/description/{code}", ctrl.GetInviteByCode,
		option.Summary("Get secret friend details by its invite-code"),
		groupTag, web.OptionAuthToken(),
	)

	fuego.Get(
		server, "/{id}", ctrl.GetSecretFriend,
		option.Summary("Get secret friend details"),
		option.Description("Get details of a specific secret friend event"),
		optionEventID, groupTag, web.OptionAuthToken(),
	)

	fuego.Patch(
		server, "/{id}", ctrl.UpdateSecretFriend,
		option.Summary("Update secret friend event"),
		option.Description("Update details of a secret friend event"),
		optionEventID, groupTag, web.OptionAuthToken(),
	)

	fuego.Post(
		server, "/{id}/draw", ctrl.DrawSecretFriend,
		option.Summary("Execute secret friend draw"),
		option.Description("Execute the draw algorithm for a secret friend event"),
		optionEventID, groupTag, web.OptionAuthToken(),
	)

	fuego.Get(
		server, "/{id}/draw-result", ctrl.GetDrawResult,
		option.Summary("Get draw result"),
		option.Description("Get the draw result for the authenticated user"),
		optionEventID, groupTag, web.OptionAuthToken(),
	)

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

func (r *Router) SubRouters() (subgroupPattern string, routers []web.RouterContract) {
	return "/{id}", []web.RouterContract{
		wishlistctrl.NewRouter(r.injector),
		denylistctrl.NewRouter(r.injector),
		participantsctrl.NewRouter(r.injector),
	}
}
