package stdrunners

import (
	"errors"
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
)

// LoginRunner returns a standard login runner for auth tests.
func LoginRunner(baseURL, email, password string, opts ...netoche.Option) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(
			http.MethodPost,
			"/auth/login",
			controllers.FormUser{
				Email:    email,
				Password: password,
			},
		),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(
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
	return netoche.New(baseURL, baseOpts...)
}
