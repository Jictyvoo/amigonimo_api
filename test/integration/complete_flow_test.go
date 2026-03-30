package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/authctrl/controllers"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/participantsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/secretfriendsctrl"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	participantsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/participants"
	secretfriendsrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/secretfriends"
)

// TestCompleteSecretFriendFlow exercises the full user journey without any DB seeding:
//  1. Signup manager + 2 participants
//  2. Manager logs in and creates a secret friend event
//  3. Participants join the event via the invite code returned by Create
//  4. One participant adds a wishlist item
//  5. Manager executes the draw
//  6. Manager and participant each retrieve their draw result
func TestCompleteSecretFriendFlow(t *testing.T) {
	engine := NewEngine(t)

	const (
		managerEmail  = "flow-manager@example.com"
		managerPass   = "flow-manager-pass"
		part1Email    = "flow-part1@example.com"
		part1Pass     = "flow-part1-pass"
		part2Email    = "flow-part2@example.com"
		part2Pass     = "flow-part2-pass"
		eventName     = "Complete Flow Event"
		wishlistLabel = "A nice book"
	)

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			// ── Step 1: signup all users ───────────────────────────────────────────
			authrunner.SignUp(
				engine.BaseURL(), controllers.FormUser{
					Email:    managerEmail,
					Username: "flow-manager",
					Password: managerPass,
				},
			),
			authrunner.SignUp(
				engine.BaseURL(), controllers.FormUser{
					Email:    part1Email,
					Username: "flow-part1",
					Password: part1Pass,
				},
			),
			authrunner.SignUp(
				engine.BaseURL(), controllers.FormUser{
					Email:    part2Email,
					Username: "flow-part2",
					Password: part2Pass,
				},
			),

			// ── Step 2: manager creates the event and joins it ────────────────────
			authrunner.Login(engine.BaseURL(), managerEmail, managerPass),
			secretfriendsrunner.Create(
				engine.BaseURL(),
				secretfriendsctrl.CreateSecretFriendRequest{
					Name:     eventName,
					Datetime: time.Now().Add(72 * time.Hour),
					Location: "Online",
				},
				netoche.ExpectStatus(http.StatusOK),
			),
			// Manager must also confirm participation to count towards the draw minimum
			participantsrunner.Confirm(
				engine.BaseURL(),
				participantsctrl.ConfirmParticipationRequest{Confirm: true},
				netoche.WithPathParamFromCtx(
					"secretFriendId",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
			),

			// ── Step 3a: participant 1 joins ───────────────────────────────────────
			authrunner.Login(engine.BaseURL(), part1Email, part1Pass),
			participantsrunner.Confirm(
				engine.BaseURL(),
				participantsctrl.ConfirmParticipationRequest{Confirm: true},
				netoche.WithPathParamFromCtx(
					"secretFriendId",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodPost,
					"/secret-friends/{id}/wishlist/",
					wishlistctrl.WishlistItemRequest{
						Label:    wishlistLabel,
						Comments: "Paperback preferred",
					},
				),
				authrunner.WithAuthHeaderFromLogin(),
				netoche.WithPathParamFromCtx(
					"id",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
				netoche.ExpectStatus(http.StatusOK),
			),

			// ── Step 3b: participant 2 joins ───────────────────────────────────────
			authrunner.Login(engine.BaseURL(), part2Email, part2Pass),
			participantsrunner.Confirm(
				engine.BaseURL(),
				participantsctrl.ConfirmParticipationRequest{Confirm: true},
				netoche.WithPathParamFromCtx(
					"secretFriendId",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
			),

			// ── Step 5: manager executes the draw ──────────────────────────────────
			authrunner.Login(engine.BaseURL(), managerEmail, managerPass),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/draw", struct{}{}),
				authrunner.WithAuthHeaderFromLogin(),
				netoche.WithPathParamFromCtx(
					"id",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
				netoche.ExpectStatus(http.StatusOK),
			),

			// ── Step 6a: manager retrieves draw result ─────────────────────────────
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/draw-result", struct{}{}),
				authrunner.WithAuthHeaderFromLogin(),
				netoche.WithPathParamFromCtx(
					"id",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
				netoche.ExpectStatus(http.StatusOK),
			),

			// ── Step 6b: participant 1 retrieves their draw result ─────────────────
			authrunner.Login(engine.BaseURL(), part1Email, part1Pass),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/draw-result", struct{}{}),
				authrunner.WithAuthHeaderFromLogin(),
				netoche.WithPathParamFromCtx(
					"id",
					func(r secretfriendsctrl.CreateSecretFriendResponse) string {
						return r.SecretFriendID
					},
				),
				netoche.ExpectStatus(http.StatusOK),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("Complete flow failed: %v", err)
	}
}
