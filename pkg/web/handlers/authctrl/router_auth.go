package authctrl

import (
	"fmt"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
)

type RouterAuth struct {
	config      Config
	middlewares []web.HttpMiddleware
}

func NewAuthRouter(config Config) *RouterAuth {
	remy.RegisterConstructorArgs2(
		config.Injector,
		remy.LazySingleton[controllers.AuthenticationController],
		controllers.NewAuthController,
	)
	return &RouterAuth{middlewares: []web.HttpMiddleware{}, config: config}
}

func (r *RouterAuth) RegisterRoutes(server *fuego.Server) error {
	authHandlers, err := remy.Get[controllers.AuthenticationController](r.config.Injector)
	if err != nil {
		return fmt.Errorf("register auth handler: %w", err)
	}

	tag := option.Tags("Authentication")
	r.registerBasicAuthRoutes(server, authHandlers, tag)
	r.registerPasswordRoutes(server, authHandlers, tag)
	r.registerTokenRoutes(server, authHandlers, tag)
	r.registerVerificationRoutes(server, authHandlers, tag)

	return nil
}

func (r *RouterAuth) GroupName() string {
	return "/auth"
}

func (r *RouterAuth) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}

func (r *RouterAuth) registerBasicAuthRoutes(
	server *fuego.Server, handlers controllers.AuthenticationController, tag fuego.RouteOption,
) {
	if r.config.ActiveRoutes.is(RouteSignUp) {
		fuego.Post(
			server,
			"/sign",
			handlers.SignUp,
			option.Summary("Starts a client signup"),
			option.Description(
				"Receive a client form, then try to decode it and create on database",
			),
			tag,
		)
	}
	if r.config.ActiveRoutes.is(RouteLogin) {
		fuego.Post(
			server, "/login", handlers.LogIn,
			option.Summary("Realizes user login"),
			option.Description(
				"Receive a client form, and then realize a login that will return a token on Authorization header",
			),
			tag,
		)
	}
}

func (r *RouterAuth) registerPasswordRoutes(
	server *fuego.Server, handlers controllers.AuthenticationController, tag fuego.RouteOption,
) {
	if r.config.ActiveRoutes.is(RouteForgotPassword) {
		fuego.Put(
			server, "/password/forgot", handlers.ForgotPassword,
			option.Summary("Initiate password recovery process"),
			option.Description(
				"Initiate the password recovery process by providing a username or email",
			),
			tag,
		)
	}
	if r.config.ActiveRoutes.is(RouteResetPassword) {
		fuego.Put(
			server, "/password/reset", handlers.PasswordReset,
			option.Summary("Reset user's password"),
			option.Description("Reset a user's password by providing a reset request"),
			tag,
		)
		fuego.Patch(
			server, "/password/check-recovery", handlers.CheckRecoveryCode,
			option.Summary("Check if given recovery code is valid"),
			option.Description(
				"Takes a recovery code and checks if it is a valid one for the given identifier",
			),
			tag,
		)
	}
}

func (r *RouterAuth) registerTokenRoutes(
	server *fuego.Server, handlers controllers.AuthenticationController, tag fuego.RouteOption,
) {
	if r.config.ActiveRoutes.is(RouteRegenerateToken) {
		fuego.Patch(
			server, "/regenerate", handlers.RegenerateAuthToken,
			option.Summary("Regenerate authentication token"),
			option.Description("Regenerate an authentication token using a refresh token"),
			tag, web.OptionAuthToken(),
		)
	}
}

func (r *RouterAuth) registerVerificationRoutes(
	server *fuego.Server, handlers controllers.AuthenticationController, tag fuego.RouteOption,
) {
	if r.config.ActiveRoutes.is(RouteVerifyUser) {
		fuego.Get(
			server, "/verify/:verify_code", handlers.VerifyUser,
			option.Summary("Verify user with a verification code"),
			option.Description("Verify a user by providing a verification code"),
			option.Path("verify_code", "Base64-encoded verification code", param.Required()),
			tag,
		)
	}
}
