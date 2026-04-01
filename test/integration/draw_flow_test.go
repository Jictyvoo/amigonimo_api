package integration

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	secretfriendsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/secretfriends"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

// TestDrawRoutesRequireAuth verifies that the draw and draw-result endpoints
// return 401 Unauthorized when no Authorization header is provided.
func TestDrawRoutesRequireAuth(t *testing.T) {
	engine := NewEngine(t)
	someID, _ := entities.NewHexID()

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/draw", struct{}{}),
				netoche.WithPathParam("id", someID),
				netoche.ExpectStatus(http.StatusUnauthorized),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/draw-result", struct{}{}),
				netoche.WithPathParam("id", someID),
				netoche.ExpectStatus(http.StatusUnauthorized),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func seedDrawEvent(
	t *testing.T,
	engine interface {
		Seed(items ...any) error
	},
	owner, userTwo, userThree *fixturesets.User,
	eventName string,
) (*fixturesets.User, *fixturesets.User, *fixturesets.User, entities.HexID) {
	t.Helper()

	drawEvent := fixtures.NewSecretFriend().
		WithOwner(owner.User).
		WithName(eventName).
		Build()
	drawEvent.Status = string(entities.StatusOpen)
	drawEvent.Datetime = sql.NullTime{Time: time.Now().Add(24 * time.Hour), Valid: true}

	ownerParticipant := fixtures.NewParticipant().
		WithUser(owner.User).
		WithSecretFriend(drawEvent).
		Build()
	p2Participant := fixtures.NewParticipant().
		WithUser(userTwo.User).
		WithSecretFriend(drawEvent).
		Build()
	p3Participant := fixtures.NewParticipant().
		WithUser(userThree.User).
		WithSecretFriend(drawEvent).
		Build()

	if err := engine.Seed(
		owner.User, owner.Profile,
		userTwo.User, userTwo.Profile,
		userThree.User, userThree.Profile,
		drawEvent,
		ownerParticipant, p2Participant, p3Participant,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	eventID, _ := entities.NewHexIDFromBytes(drawEvent.ID)
	return owner, userTwo, userThree, eventID
}

func TestDrawSucceedsAndResultRetrievable(t *testing.T) {
	engine := NewEngine(t)
	const password = "draw-success-password"

	owner := fixturesets.NewUser("draw-success-owner@example.com", password, "Draw Owner")
	p2 := fixturesets.NewUser("draw-success-p2@example.com", password, "Participant Two")
	p3 := fixturesets.NewUser("draw-success-p3@example.com", password, "Participant Three")

	_, _, _, drawEventID := seedDrawEvent(t, engine, owner, p2, p3, "Draw Success Event")

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, password),
			secretfriendsrunner.Draw(
				engine.BaseURL(),
				drawEventID,
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					secretfriendsctrl.DrawSecretFriendResponse{
						SecretFriendID: drawEventID.String(),
						Status:         string(entities.StatusDrawn),
						ResultCount:    3,
					},
				),
			),
			// After draw, owner can retrieve their result
			secretfriendsrunner.GetDrawResult(
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

func TestDrawAlreadyDrawnConflict(t *testing.T) {
	engine := NewEngine(t)
	const password = "draw-conflict-password"

	owner := fixturesets.NewUser("draw-conflict-owner@example.com", password, "Conflict Owner")
	p2 := fixturesets.NewUser("draw-conflict-p2@example.com", password, "Conflict P2")
	p3 := fixturesets.NewUser("draw-conflict-p3@example.com", password, "Conflict P3")

	_, _, _, drawEventID := seedDrawEvent(t, engine, owner, p2, p3, "Draw Conflict Event")

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, password),
			// First draw succeeds
			secretfriendsrunner.Draw(
				engine.BaseURL(),
				drawEventID,
				netoche.ExpectStatus(http.StatusOK),
			),
			// Second draw on an already-drawn event must be rejected
			secretfriendsrunner.Draw(
				engine.BaseURL(),
				drawEventID,
				netoche.ExpectStatus(http.StatusConflict),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

// TestDrawResultBeforeDraw verifies that requesting a draw result before the draw
// has been executed returns 404 Not Found.
func TestDrawResultBeforeDraw(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "draw-before-password"

	owner := fixturesets.NewUser("draw-before-owner@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(
		owner,
		fixturesets.NewUser("draw-before-p1@example.com", userPassword, ""),
		"Draw Before Event",
	)

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), owner.User.Email, userPassword),
			secretfriendsrunner.GetDrawResult(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				netoche.ExpectStatus(http.StatusNotFound),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
