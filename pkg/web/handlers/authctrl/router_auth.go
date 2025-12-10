package authctrl

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"

	"github.com/jictyvoo/amigonimo_api/pkg/web"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
)

type RouterAuth struct {
	config      Config
	middlewares []web.HttpMiddleware
}

func NewAuthRouter(config Config) *RouterAuth {
	return &RouterAuth{middlewares: []web.HttpMiddleware{}, config: config}
}

func (r *RouterAuth) RegisterRoutes(server *fuego.Server) error {
	authHandlers := controllers.NewAuthController(r.config.SecretKey, nil, nil, nil)
	groupTag := option.Tags("Authentication")
	if r.config.ActiveRoutes.is(RouteSignUp) {
		fuego.Post(
			server, "/sign", authHandlers.SignUp,
			option.Summary("Starts a client signup"),
			option.Description(
				"Receive a client form, then try to decode it and create on database",
			),
			groupTag,
		)
	}
	if r.config.ActiveRoutes.is(RouteLogin) {
		fuego.Post(
			server, "/login", authHandlers.LogIn,
			option.Summary("Realizes user login"),
			option.Description(
				"Receive a client form, and then realize a login that will return a token on Authorization header",
			),
			groupTag,
		)
	}

	if r.config.ActiveRoutes.is(RouteForgotPassword) {
		fuego.Put(
			server, "/password/forgot", authHandlers.ForgotPassword,
			option.Summary("Initiate password recovery process"),
			option.Description(
				"Initiate the password recovery process by providing a username or email",
			),
			groupTag,
		)
	}
	if r.config.ActiveRoutes.is(RouteResetPassword) {
		fuego.Put(
			server, "/password/reset", authHandlers.PasswordReset,
			option.Summary("Reset user's password"),
			option.Description("Reset a user's password by providing a reset request"),
			groupTag,
		)
		fuego.Patch(
			server, "/password/check-recovery", authHandlers.CheckRecoveryCode,
			option.Summary("Check if given recovery code is valid"),
			option.Description(
				"Takes a recovery code and checks if it is a valid one for the given identifier",
			),
			groupTag,
		)
	}
	if r.config.ActiveRoutes.is(RouteRegenerateToken) {
		fuego.Patch(
			server, "/regenerate", authHandlers.RegenerateAuthToken,
			option.Summary("Regenerate authentication token"),
			option.Description("Regenerate an authentication token using a refresh token"),
			groupTag, web.OptionAuthToken(),
		)
	}
	if r.config.ActiveRoutes.is(RouteVerifyUser) {
		fuego.Get(
			server, "/verify/:verify_code", authHandlers.VerifyUser,
			option.Summary("Verify user with a verification code"),
			option.Description("Verify a user by providing a verification code"),
			option.Path("verify_code", "Base64-encoded verification code", param.Required()),
			groupTag,
		)
	}

	return nil
}

func (r *RouterAuth) GroupName() string {
	return "/auth"
}

func (r *RouterAuth) Middlewares() []web.HttpMiddleware {
	return r.middlewares
}
