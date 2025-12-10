package authctrl

import "time"

type DefinedRoute uint16

const (
	RouteLogin DefinedRoute = 1 << iota
	RouteSignUp
	RouteForgotPassword
	RouteResetPassword
	RouteEditPassword
	RouteEditEmail
	RouteEditUsername
	RouteRegenerateToken
	RouteVerifyUser
)

func (r DefinedRoute) is(other DefinedRoute) bool {
	return r&other == other
}

type (
	Config struct {
		ActiveRoutes DefinedRoute
		SecretKey    []byte
	}
)

func TokenInitialExpiration() time.Time {
	return time.Date(1970, time.January, 1, 1, 0, 0, 0, time.UTC)
}
