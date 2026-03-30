package usersrunner

import (
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
)

func EditPassword(
	baseURL string,
	req controllers.FormEditPassword,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPatch, "/users/edit/password", req),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(controllers.SuccessResponse{
			Success: true,
			Message: "Password changed successfully",
		}),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func EditUsername(
	baseURL string,
	req controllers.FormEditUsername,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPatch, "/users/edit/username", req),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(controllers.SuccessResponse{
			Success: true,
			Message: "Username changed successfully",
		}),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func EditEmail(
	baseURL string,
	req controllers.FormEditEmail,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPatch, "/users/edit/email", req),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(controllers.SuccessResponse{
			Success: true,
			Message: "Email changed successfully",
		}),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}
