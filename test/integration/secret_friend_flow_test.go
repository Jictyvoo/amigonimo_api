package integration

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

func TestCreateSecretFriendAndJoin(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "generic@test.password"

	// Users
	manager := fixtures.NewUser().
		WithUsername("manager").
		WithEmail("manager@example.com").
		WithPassword(userPassword).
		Build()
	participant := fixtures.NewUser().
		WithUsername("participant").
		WithEmail("participant@example.com").
		WithPassword(userPassword).
		Build()
	engine.Seed(manager, participant)

	mr := runners.MultiRunner{
		Runners: []runners.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), manager.Email, userPassword),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(
					http.MethodPost,
					"/secret-friends/",
					secretfriendsctrl.CreateSecretFriendRequest{
						Name:     "Integration Test Event",
						Datetime: time.Now().Add(24 * time.Hour),
						Location: "Virtual",
					},
				),
				reqrunner.WithHeaderFromCtx(
					"Authorization",
					func(logResp controllers.LoginResponse) string { return "Bearer " + logResp.Token },
				),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
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
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(
					http.MethodPost,
					"/secret-friends/{secretFriendId}/participants/",
					participantsctrl.ConfirmParticipationRequest{Confirm: true},
				),
				reqrunner.WithHeaderFromCtx(
					"Authorization",
					func(logResp controllers.LoginResponse) string { return "Bearer " + logResp.Token },
				),
				reqrunner.WithPathParamFromCtx(
					"secretFriendId",
					func(newSecretFriend secretfriendsctrl.CreateSecretFriendResponse) string {
						return newSecretFriend.SecretFriendID
					},
				),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					participantsctrl.ConfirmParticipationResponse{Success: true},
					func(expected, actual *participantsctrl.ConfirmParticipationResponse) error {
						expected.ParticipantID = actual.ParticipantID
						return nil
					},
				),
			),
		},
	}

	if err := mr.Run(t); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
