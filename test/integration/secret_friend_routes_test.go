package integration

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	participantsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/participants"
	secretfriendsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/secretfriends"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

func TestSecretFriendMountedRoutes(t *testing.T) {
	engine := NewEngine(t)
	const ownerPassword = "owner-routes-password"

	owner := fixturesets.NewUser("routes-owner@example.com", ownerPassword, "Routes Owner")
	target := fixturesets.NewUser("routes-target@example.com", ownerPassword, "Target Person")
	other := fixturesets.NewUser("routes-other@example.com", ownerPassword, "Other Person")
	routeSet := fixturesets.NewSecretFriendRoutes(owner, target, other)

	if err := engine.Seed(routeSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	openEventID := routeSet.OpenEventID()
	ownerID := routeSet.Owner.ID()
	targetUserID := routeSet.Target.ID()
	otherUserID := routeSet.Other.ID()

	updatedName := "Updated Open Routes Event"
	updatedLocation := "Updated Place"

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), routeSet.Owner.User.Email, ownerPassword),
			secretfriendsrunner.InviteInfo(
				engine.BaseURL(),
				routeSet.OpenEvent.InviteCode,
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.InviteInfoResponse{
						SecretFriendID: openEventID.String(),
						Name:           routeSet.OpenEvent.Name,
					},
				),
			),
			secretfriendsrunner.Get(
				engine.BaseURL(),
				openEventID,
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.GetSecretFriendResponse{
						ID:                openEventID.String(),
						Name:              routeSet.OpenEvent.Name,
						Location:          routeSet.OpenEvent.Location.String,
						OwnerID:           ownerID.String(),
						ParticipantsCount: 3,
						Status:            routeSet.OpenEvent.Status,
					},
					func(expected, actual *secretfriendsctrl.GetSecretFriendResponse) error {
						expected.Datetime = actual.Datetime
						return nil
					},
				),
			),
			secretfriendsrunner.Update(
				engine.BaseURL(),
				openEventID,
				secretfriendsctrl.UpdateSecretFriendRequest{
					Name:     updatedName,
					Location: updatedLocation,
				},
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.UpdateSecretFriendResponse{
						Success: true,
						Message: "secret friend updated successfully",
					},
				),
			),
			secretfriendsrunner.Get(
				engine.BaseURL(),
				openEventID,
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.GetSecretFriendResponse{
						ID:                openEventID.String(),
						Name:              updatedName,
						Location:          updatedLocation,
						OwnerID:           ownerID.String(),
						ParticipantsCount: 3,
						Status:            routeSet.OpenEvent.Status,
					},
					func(expected, actual *secretfriendsctrl.GetSecretFriendResponse) error {
						expected.Datetime = actual.Datetime
						return nil
					},
				),
			),
			participantsrunner.List(
				engine.BaseURL(),
				openEventID,
				netoche.ExpectStatus(http.StatusOK),
				participantsrunner.ExpectList(
					[]participantsctrl.ParticipantResponse{
						{
							ParticipantID: routeSet.OpenOwnerEntryID().String(),
							UserID:        ownerID.String(),
							Fullname:      routeSet.Owner.Profile.Fullname.String,
						},
						{
							ParticipantID: routeSet.OpenTargetEntryID().String(),
							UserID:        targetUserID.String(),
							Fullname:      routeSet.Target.Profile.Fullname.String,
						},
						{
							ParticipantID: routeSet.OpenOtherEntryID().String(),
							UserID:        otherUserID.String(),
							Fullname:      routeSet.Other.Profile.Fullname.String,
						},
					},
				),
			),
			// Invalid UUID returns 400
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}", struct{}{}),
				authrunner.WithAuthHeaderFromLogin(),
				netoche.WithPathParam("id", "not-a-uuid"),
				netoche.ExpectStatus(http.StatusBadRequest),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestListSecretFriends(t *testing.T) {
	engine := NewEngine(t)
	const ownerPassword = "list-events-password"

	owner := fixturesets.NewUser("list-owner@example.com", ownerPassword, "List Owner")
	other := fixturesets.NewUser("list-other@example.com", ownerPassword, "Other User")
	eventSet := fixturesets.NewOwnerParticipant(owner, other, "Listed Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, ownerPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/", struct{}{}),
				authrunner.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusOK),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestDrawSecretFriendRoute(t *testing.T) {
	engine := NewEngine(t)
	const ownerPassword = "draw-routes-password"

	owner := fixturesets.NewUser("draw-owner@example.com", ownerPassword, "Draw Owner")
	userTwo := fixturesets.NewUser("draw-user-two@example.com", ownerPassword, "User Two")
	userThree := fixturesets.NewUser("draw-user-three@example.com", ownerPassword, "User Three")

	drawEvent := fixtures.NewSecretFriend().
		WithOwner(owner.User).
		WithName("Draw Execution Event").
		Build()
	drawEvent.Status = string(entities.StatusOpen)
	drawEvent.Datetime = sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true}

	ownerParticipant := fixtures.NewParticipant().
		WithUser(owner.User).
		WithSecretFriend(drawEvent).
		Build()
	userTwoParticipant := fixtures.NewParticipant().
		WithUser(userTwo.User).
		WithSecretFriend(drawEvent).
		Build()
	userThreeParticipant := fixtures.NewParticipant().
		WithUser(userThree.User).
		WithSecretFriend(drawEvent).
		Build()

	if err := engine.Seed(
		owner.User,
		owner.Profile,
		userTwo.User,
		userTwo.Profile,
		userThree.User,
		userThree.Profile,
		drawEvent,
		ownerParticipant,
		userTwoParticipant,
		userThreeParticipant,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	drawEventID, _ := entities.NewHexIDFromBytes(drawEvent.ID)

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, ownerPassword),
			secretfriendsrunner.Draw(
				engine.BaseURL(),
				drawEventID,
				netoche.ExpectStatus(http.StatusOK),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestGetSecretFriendNotFound(t *testing.T) {
	engine := NewEngine(t)
	const ownerPassword = "notfound-password"

	user := fixturesets.NewUser("notfound-user@example.com", ownerPassword, "")
	if err := engine.Seed(user.User, user.Profile); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	nonExistentID, _ := entities.NewHexID()

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), user.User.Email, ownerPassword),
			secretfriendsrunner.Get(
				engine.BaseURL(),
				nonExistentID,
				netoche.ExpectStatus(http.StatusNotFound),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestUpdateSecretFriendByNonOwner(t *testing.T) {
	engine := NewEngine(t)
	const password = "nonowner-update-password"

	owner := fixturesets.NewUser("nonowner-event-owner@example.com", password, "Event Owner")
	participant := fixturesets.NewUser("nonowner-participant@example.com", password, "Participant")
	eventSet := fixturesets.NewOwnerParticipant(owner, participant, "Non-owner Update Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), participant.User.Email, password),
			secretfriendsrunner.Update(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				secretfriendsctrl.UpdateSecretFriendRequest{Name: "Hijacked Name"},
				netoche.ExpectStatus(http.StatusForbidden),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
