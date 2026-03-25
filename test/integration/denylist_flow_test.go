package integration

import (
	"net/http"
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func TestDenylistFlowSeeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "denylist-flow-password"

	ownerBuilder := fixtures.NewUser().
		WithEmail("denylist-owner@example.com").
		WithPassword(userPassword)
	owner := ownerBuilder.Build()
	ownerProfile := ownerBuilder.BuildProfile()
	participantBuilder := fixtures.NewUser().
		WithEmail("denylist-participant@example.com").
		WithPassword(userPassword)
	participant := participantBuilder.Build()
	participantProfile := participantBuilder.BuildProfile()
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

	if err := engine.Seed(
		owner,
		ownerProfile,
		participant,
		participantProfile,
		secretFriend,
		ownerEntry,
		participantEntry,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	secretFriendID, _ := entities.NewHexIDFromBytes(secretFriend.ID)
	deniedUserID, _ := entities.NewHexIDFromBytes(participant.ID)

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), owner.Email, userPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPost,
					"/secret-friends/{id}/denylist/",
					denylistctrl.AddDenyListRequest{TargetUserID: deniedUserID.String()},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", secretFriendID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					denylistctrl.DeniedUserResponse{
						UserID: deniedUserID.String(),
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/denylist/", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", secretFriendID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					[]denylistctrl.DeniedUserResponse{
						{
							UserID:   deniedUserID.String(),
							Fullname: participantProfile.Fullname.String,
						},
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodDelete,
					"/secret-friends/{id}/denylist/{deniedUserId}",
					struct{}{},
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", secretFriendID),
				netoche.WithPathParam("deniedUserId", deniedUserID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					denylistctrl.RemoveDenyListEntryResponse{
						Success:   true,
						DeletedID: deniedUserID.String(),
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/denylist/", struct{}{}),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", secretFriendID),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					[]denylistctrl.DeniedUserResponse{},
				),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
