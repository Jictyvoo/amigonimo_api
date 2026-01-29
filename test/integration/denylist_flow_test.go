package integration

import (
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

func TestDenylistFlowSeeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "denylist-flow-password"

	owner := fixtures.NewUser().
		WithEmail("denylist-owner@example.com").
		WithPassword(userPassword).
		Build()
	participant := fixtures.NewUser().
		WithEmail("denylist-participant@example.com").
		WithPassword(userPassword).
		Build()
	secretFriend := fixtures.NewSecretFriend().
		WithOwner(owner).
		WithName("Denylist Flow Event").
		Build()
	ownerEntry := fixtures.NewParticipant().
		WithUser(owner).
		WithSecretFriend(secretFriend).
		Build()
	participantEntry := fixtures.NewParticipant().
		WithUser(participant).
		WithSecretFriend(secretFriend).
		Build()

	engine.Seed(owner, participant, secretFriend, ownerEntry, participantEntry)

	secretFriendID, _ := entities.NewHexIDFromBytes(secretFriend.ID)
	deniedUserID, _ := entities.NewHexIDFromBytes(participant.ID)

	mr := runners.MultiRunner{
		Runners: []runners.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), owner.Email, userPassword),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(
					http.MethodPost,
					"/secret-friends/{id}/denylist/",
					denylistctrl.AddDenyListRequest{TargetUserID: deniedUserID.String()},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.WithPathParam("id", secretFriendID),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					denylistctrl.DeniedUserResponse{
						UserID: deniedUserID.String(),
						// Fullname: participant.Fullname, // TODO: Analyze if is really expected to return more info
					},
				),
			),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(http.MethodGet, "/secret-friends/{id}/denylist/", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.WithPathParam("id", secretFriendID),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					[]denylistctrl.DeniedUserResponse{
						{
							UserID:   deniedUserID.String(),
							Fullname: participant.Fullname,
						},
					},
				),
			),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(
					http.MethodDelete,
					"/secret-friends/{id}/denylist/{deniedUserId}",
					struct{}{},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.WithPathParam("id", secretFriendID),
				reqrunner.WithPathParam("deniedUserId", deniedUserID),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					denylistctrl.RemoveDenyListEntryResponse{
						Success:   true,
						DeletedID: deniedUserID.String(),
					},
				),
			),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(http.MethodGet, "/secret-friends/{id}/denylist/", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.WithPathParam("id", secretFriendID),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					[]denylistctrl.DeniedUserResponse{},
				),
			),
		},
	}

	if err := mr.Run(t); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
