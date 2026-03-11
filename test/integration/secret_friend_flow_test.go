package integration

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func TestCreateSecretFriendAndJoin(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "generic@test.password"

	// Users
	managerBuilder := fixtures.NewUser().
		WithUsername("manager").
		WithEmail("manager@example.com").
		WithPassword(userPassword)
	manager := managerBuilder.Build()
	managerProfile := managerBuilder.BuildProfile()
	participantBuilder := fixtures.NewUser().
		WithUsername("participant").
		WithEmail("participant@example.com").
		WithPassword(userPassword)
	participant := participantBuilder.Build()
	participantProfile := participantBuilder.BuildProfile()
	engine.Seed(manager, managerProfile, participant, participantProfile)

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), manager.Email, userPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPost,
					"/secret-friends/",
					secretfriendsctrl.CreateSecretFriendRequest{
						Name:     "Integration Test Event",
						Datetime: time.Now().Add(24 * time.Hour),
						Location: "Virtual",
					},
				),
				stdrunners.WithAuthHeaderFromLogin(),
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
			),
			stdrunners.LoginRunner(engine.BaseURL(), participant.Email, userPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPost,
					"/secret-friends/{secretFriendId}/participants/",
					participantsctrl.ConfirmParticipationRequest{Confirm: true},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParamFromCtx(
					"secretFriendId",
					func(newSecretFriend secretfriendsctrl.CreateSecretFriendResponse) string {
						return newSecretFriend.SecretFriendID
					},
				),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					participantsctrl.ConfirmParticipationResponse{Success: true},
					func(expected, actual *participantsctrl.ConfirmParticipationResponse) error {
						expected.ParticipantID = actual.ParticipantID
						return nil
					},
				),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
