package stdrunners

import (
	"errors"
	"net/http"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

// LoginRunner returns a standard login runner for auth tests.
func LoginRunner(baseURL, email, password string, opts ...reqrunner.Option) *reqrunner.HttpRunner {
	baseOpts := []reqrunner.Option{
		reqrunner.WithRequest(
			http.MethodPost,
			"/auth/login",
			controllers.FormUser{
				Email:    email,
				Password: password,
			},
		),
		reqrunner.ExpectStatus(http.StatusOK),
		reqrunner.ExpectBody(
			controllers.LoginResponse{},
			func(expected, actual *controllers.LoginResponse) error {
				if actual.Token == "" {
					return errors.New("token is empty")
				}

				actual.Token = expected.Token
				return nil
			},
		),
	}

	baseOpts = append(baseOpts, opts...)
	return reqrunner.NewHttpRunner(baseURL, baseOpts...)
}
