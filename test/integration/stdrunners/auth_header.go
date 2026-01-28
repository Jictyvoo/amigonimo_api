package stdrunners

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

func WithAuthHeaderFromLogin() reqrunner.Option {
	return reqrunner.WithHeaderFromCtx(
		"Authorization",
		func(loginResp controllers.LoginResponse) string {
			return fmt.Sprintf("Bearer %s", loginResp.Token)
		},
	)
}
