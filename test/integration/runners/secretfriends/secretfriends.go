package secretfriendsrunner

import (
	"errors"
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
)

func Create(
	baseURL string,
	req secretfriendsctrl.CreateSecretFriendRequest,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPost, "/secret-friends/", req),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(
			secretfriendsctrl.CreateSecretFriendResponse{},
			func(expected, actual *secretfriendsctrl.CreateSecretFriendResponse) error {
				if actual.SecretFriendID == "" {
					return errors.New("secret friend id is empty")
				}
				if actual.InviteCode == "" {
					return errors.New("invite code is empty")
				}
				*expected = *actual
				return nil
			},
		),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func Get(baseURL string, id entities.HexID, opts ...netoche.Option) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodGet, "/secret-friends/{id}", struct{}{}),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.WithPathParam("id", id),
	}
	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func Update(
	baseURL string,
	id entities.HexID,
	req secretfriendsctrl.UpdateSecretFriendRequest,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPatch, "/secret-friends/{id}", req),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.WithPathParam("id", id),
	}
	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func InviteInfo(baseURL, inviteCode string, opts ...netoche.Option) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(
			http.MethodGet,
			"/secret-friends/invites/description/{code}",
			struct{}{},
		),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.WithPathParam("code", inviteCode),
	}
	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func GetDrawResult(
	baseURL string,
	id entities.HexID,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/draw-result", struct{}{}),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.WithPathParam("id", id),
	}
	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func Draw(baseURL string, id entities.HexID, opts ...netoche.Option) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/draw", struct{}{}),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.WithPathParam("id", id),
	}
	return netoche.New(baseURL, append(baseOpts, opts...)...)
}
