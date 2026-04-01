package integration

import (
	"net/http"
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/denylistctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	denylistrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/denylist"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

// TestDenylistRequiresAuth verifies that denylist endpoints return 401 Unauthorized
// when no Authorization header is provided.
func TestDenylistRequiresAuth(t *testing.T) {
	engine := NewEngine(t)
	someID, _ := entities.NewHexID()

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/denylist/", struct{}{}),
				netoche.WithPathParam("id", someID),
				netoche.ExpectStatus(http.StatusUnauthorized),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/denylist/", struct{}{}),
				netoche.WithPathParam("id", someID),
				netoche.ExpectStatus(http.StatusUnauthorized),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestDenylistFlowSeeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "denylist-flow-password"

	owner := fixturesets.NewUser("denylist-owner@example.com", userPassword, "")
	participant := fixturesets.NewUser("denylist-participant@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(owner, participant, "Denylist Flow Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, userPassword),
			denylistrunner.AddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: participant.ID().String()},
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					denylistctrl.DeniedUserResponse{
						UserID: participant.ID().String(),
					},
				),
			),
			denylistrunner.ListEntries(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				netoche.ExpectStatus(http.StatusOK),
				denylistrunner.ExpectEntries(
					[]denylistctrl.DeniedUserResponse{
						{
							UserID:   participant.ID().String(),
							Fullname: participant.Profile.Fullname.String,
						},
					},
				),
			),
			denylistrunner.RemoveEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				participant.ID(),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					denylistctrl.RemoveDenyListEntryResponse{
						Success:   true,
						DeletedID: participant.ID().String(),
					},
				),
			),
			denylistrunner.ListEntries(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				netoche.ExpectStatus(http.StatusOK),
				denylistrunner.ExpectEntries([]denylistctrl.DeniedUserResponse{}),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestDenylistSelfEntry(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "denylist-self-password"

	owner := fixturesets.NewUser("denylist-self-owner@example.com", userPassword, "")
	other := fixturesets.NewUser("denylist-self-other@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(owner, other, "Self Denylist Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, userPassword),
			denylistrunner.FailedAddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: owner.ID().String()},
				http.StatusBadRequest,
				"you cannot add yourself to the denylist",
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

// TestDenylistCapacityExceeded seeds 4 participants with MaxDenyListSize=4.
// The 50% cap limits the effective max to 2, so the 3rd AddEntry must be rejected.
func TestDenylistCapacityExceeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "denylist-cap-password"

	owner := fixturesets.NewUser("denylist-cap-owner@example.com", userPassword, "")
	p1 := fixturesets.NewUser("denylist-cap-p1@example.com", userPassword, "")
	p2 := fixturesets.NewUser("denylist-cap-p2@example.com", userPassword, "")
	p3 := fixturesets.NewUser("denylist-cap-p3@example.com", userPassword, "")

	// 4 participants, MaxDenyListSize=4 → effective limit = 4/2 = 2
	eventSet := fixturesets.NewOwnerParticipant(owner, p1, "Cap Event").
		WithMaxDenyListSize(4)

	p2Entry := fixtures.NewParticipant().
		WithUser(p2.User).
		WithSecretFriend(eventSet.SecretFriend).
		Build()
	p3Entry := fixtures.NewParticipant().
		WithUser(p3.User).
		WithSecretFriend(eventSet.SecretFriend).
		Build()

	seedables := append(
		eventSet.Seedables(),
		p2.User, p2.Profile, p2Entry,
		p3.User, p3.Profile, p3Entry,
	)
	if err := engine.Seed(seedables...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, userPassword),
			denylistrunner.AddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: p1.ID().String()},
				netoche.ExpectStatus(http.StatusOK),
			),
			denylistrunner.AddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: p2.ID().String()},
				netoche.ExpectStatus(http.StatusOK),
			),
			// 3rd entry exceeds 50% cap (2) → must be rejected
			denylistrunner.AddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: p3.ID().String()},
				netoche.ExpectStatus(http.StatusConflict),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

// TestDenylistDuplicateEntry verifies that adding the same participant to the
// denylist twice returns 409 Conflict.
func TestDenylistDuplicateEntry(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "denylist-dup-password"

	owner := fixturesets.NewUser("denylist-dup-owner@example.com", userPassword, "")
	target := fixturesets.NewUser("denylist-dup-target@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(owner, target, "Denylist Duplicate Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, userPassword),
			// First add succeeds.
			denylistrunner.AddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: target.ID().String()},
				netoche.ExpectStatus(http.StatusOK),
			),
			// Adding the same user again must be rejected with 409 Conflict.
			denylistrunner.AddEntry(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				denylistctrl.AddDenyListRequest{TargetUserID: target.ID().String()},
				netoche.ExpectStatus(http.StatusConflict),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
