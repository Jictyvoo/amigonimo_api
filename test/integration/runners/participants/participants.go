package participantsrunner

import (
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
)

func Confirm(
	baseURL string,
	request participantsctrl.ConfirmParticipationRequest,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(
			http.MethodPost,
			"/secret-friends/{secretFriendId}/participants/",
			request,
		),
		authrunner.WithAuthHeaderFromLogin(),
		netoche.ExpectStatus(http.StatusOK),
		netoche.ExpectBody(
			participantsctrl.ConfirmParticipationResponse{Success: true},
			func(expected, actual *participantsctrl.ConfirmParticipationResponse) error {
				expected.ParticipantID = actual.ParticipantID
				return nil
			},
		),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func List(baseURL string, secretFriendID entities.HexID, opts ...netoche.Option) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/participants/", struct{}{}),
		netoche.WithPathParam("id", secretFriendID),
		authrunner.WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func ExpectList(expected []participantsctrl.ParticipantResponse) netoche.Option {
	return netoche.ExpectBody(expected)
}
