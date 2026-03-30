package integration

import (
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	participantsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/participants"
	secretfriendsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/secretfriends"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

func TestCreateSecretFriendAndJoin(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "generic@test.password"

	manager := fixturesets.NewUser("manager@example.com", userPassword, "")
	participant := fixturesets.NewUser("participant@example.com", userPassword, "")
	if err := engine.Seed(manager.User, manager.Profile, participant.User, participant.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), manager.User.Email, userPassword),
			secretfriendsrunner.Create(
				engine.BaseURL(),
				secretfriendsctrl.CreateSecretFriendRequest{
					Name:     "Integration Test Event",
					Datetime: time.Now().Add(24 * time.Hour),
					Location: "Virtual",
				},
			),
			authrunner.Login(engine.BaseURL(), participant.User.Email, userPassword),
			participantsrunner.Confirm(
				engine.BaseURL(),
				participantsctrl.ConfirmParticipationRequest{Confirm: true},
				netoche.WithPathParamFromCtx(
					"secretFriendId",
					func(newSecretFriend secretfriendsctrl.CreateSecretFriendResponse) string {
						return newSecretFriend.SecretFriendID
					},
				),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
