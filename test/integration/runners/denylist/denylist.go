package denylistrunner

import (
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func buildAddEntryRequest(
	secretFriendID entities.HexID,
	req denylistctrl.AddDenyListRequest,
) []netoche.Option {
	return []netoche.Option{
		netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/denylist/", req),
		netoche.WithPathParam("id", secretFriendID),
		authrunner.WithAuthHeaderFromLogin(),
	}
}

func AddEntry(
	baseURL string,
	secretFriendID entities.HexID,
	req denylistctrl.AddDenyListRequest,
	opts ...netoche.Option,
) atores.Runner {
	return netoche.New(baseURL, append(buildAddEntryRequest(secretFriendID, req), opts...)...)
}

// FailedAddEntry expects an error response with the given status code.
// Pass a non-empty detail to also assert the response body detail message.
func FailedAddEntry(
	baseURL string,
	secretFriendID entities.HexID,
	req denylistctrl.AddDenyListRequest,
	statusCode int,
	detail string,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := append(buildAddEntryRequest(secretFriendID, req), netoche.ExpectStatus(statusCode))
	if detail != "" {
		baseOpts = append(baseOpts, netoche.ExpectBody(fixtures.ErrorDetail{Detail: detail}))
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func ListEntries(
	baseURL string,
	secretFriendID entities.HexID,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/denylist/", struct{}{}),
		netoche.WithPathParam("id", secretFriendID),
		authrunner.WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func RemoveEntry(
	baseURL string,
	secretFriendID, deniedUserID entities.HexID,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(
			http.MethodDelete,
			"/secret-friends/{id}/denylist/{deniedUserId}",
			struct{}{},
		),
		netoche.WithPathParam("id", secretFriendID),
		netoche.WithPathParam("deniedUserId", deniedUserID),
		authrunner.WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func ExpectEntries(expected []denylistctrl.DeniedUserResponse) netoche.Option {
	return netoche.ExpectBody(expected)
}
