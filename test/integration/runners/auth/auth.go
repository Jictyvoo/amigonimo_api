package authrunner

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func WithAuthHeaderFromLogin() netoche.Option {
	return netoche.WithHeaderFromCtx(
		"Authorization",
		func(loginResp controllers.LoginResponse) string {
			return fmt.Sprintf("Bearer %s", loginResp.Token)
		},
	)
}

func buildLoginRequest(email, password string) []netoche.Option {
	return []netoche.Option{
		netoche.WithRequest(
			http.MethodPost,
			"/auth/login",
			controllers.FormUser{
				Email:    email,
				Password: password,
			},
		),
	}
}

func Login(baseURL, email, password string, opts ...netoche.Option) atores.Runner {
	baseOpts := append(
		buildLoginRequest(email, password),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(
			controllers.LoginResponse{},
			func(expected, actual *controllers.LoginResponse) error {
				if actual.Token == "" {
					return errors.New("token is empty")
				}
				expected.Token = actual.Token
				return nil
			},
		),
	)

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

// FailedLogin expects an error response with the given status code.
// Pass a non-empty detail to also assert the response body detail message.
func FailedLogin(
	baseURL, email, password string,
	statusCode int,
	detail string,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(buildLoginRequest(email, password), netoche.ExpectStatus(statusCode))
	if detail != "" {
		baseOpts = append(baseOpts, netoche.ExpectBody(fixtures.ErrorDetail{Detail: detail}))
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func buildSignUpRequest(req controllers.FormUser) []netoche.Option {
	return []netoche.Option{
		netoche.WithRequest(http.MethodPost, "/auth/sign", req),
	}
}

func SignUp(baseURL string, req controllers.FormUser, opts ...netoche.Option) atores.Runner {
	baseOpts := append(
		buildSignUpRequest(req),
		netoche.ExpectStatus(http.StatusCreated),
		netoche.ExpectBody(
			controllers.SuccessResponse{
				Success: true,
				Message: "User created successfully",
			},
		),
	)

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

// FailedSignUp expects an error response with the given status code.
// Pass a non-empty detail to also assert the response body detail message.
func FailedSignUp(
	baseURL string,
	req controllers.FormUser,
	statusCode int,
	detail string,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(buildSignUpRequest(req), netoche.ExpectStatus(statusCode))
	if detail != "" {
		baseOpts = append(baseOpts, netoche.ExpectBody(fixtures.ErrorDetail{Detail: detail}))
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func buildCheckRecoveryRequest(req controllers.FormRecoveryCode) []netoche.Option {
	return []netoche.Option{
		netoche.WithRequest(http.MethodPatch, "/auth/password/check-recovery", req),
	}
}

func CheckRecoveryCode(
	baseURL string,
	req controllers.FormRecoveryCode,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(
		buildCheckRecoveryRequest(req),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(
			controllers.SuccessResponse{
				Success: true,
				Message: "sent recovery-code is valid",
			},
		),
	)

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

// FailedCheckRecoveryCode expects an error response with the given status code.
// Pass a non-empty detail to also assert the response body detail message.
func FailedCheckRecoveryCode(
	baseURL string,
	req controllers.FormRecoveryCode,
	statusCode int,
	detail string,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(buildCheckRecoveryRequest(req), netoche.ExpectStatus(statusCode))
	if detail != "" {
		baseOpts = append(baseOpts, netoche.ExpectBody(fixtures.ErrorDetail{Detail: detail}))
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func buildResetPasswordRequest(req controllers.FormResetPassword) []netoche.Option {
	return []netoche.Option{
		netoche.WithRequest(http.MethodPut, "/auth/password/reset", req),
	}
}

func ResetPassword(
	baseURL string,
	req controllers.FormResetPassword,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(
		buildResetPasswordRequest(req),
		netoche.ExpectStatus(http.StatusCreated),
		netoche.ExpectBody(
			controllers.SuccessResponse{
				Success: true,
				Message: "password changed",
			},
		),
	)

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

// FailedResetPassword expects an error response with the given status code.
// Pass a non-empty detail to also assert the response body detail message.
func FailedResetPassword(
	baseURL string,
	req controllers.FormResetPassword,
	statusCode int,
	detail string,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(buildResetPasswordRequest(req), netoche.ExpectStatus(statusCode))
	if detail != "" {
		baseOpts = append(baseOpts, netoche.ExpectBody(fixtures.ErrorDetail{Detail: detail}))
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func ForgotPassword(
	baseURL string,
	req controllers.FormUser,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPut, "/auth/password/forgot", req),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func RegenerateToken(baseURL string, opts ...netoche.Option) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPatch, "/auth/regenerate", struct{}{}),
		WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}
