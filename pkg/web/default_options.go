package web

import (
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func OptionAuthToken() func(route *fuego.BaseRoute) {
	return option.Header(
		"Authorization", "Bearer {token}",
		param.Required(),
		param.Description("Authorization header with the authentication token"),
	)
}
