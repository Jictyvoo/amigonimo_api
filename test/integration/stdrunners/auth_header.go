package stdrunners

import (
	"fmt"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
)

func WithAuthHeaderFromLogin() netoche.Option {
	return netoche.WithHeaderFromCtx(
		"Authorization",
		func(loginResp controllers.LoginResponse) string {
			return fmt.Sprintf("Bearer %s", loginResp.Token)
		},
	)
}
