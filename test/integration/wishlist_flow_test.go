package integration

import (
	"net/http"
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	wishlistrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/wishlist"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

// TestWishlistRequiresAuth verifies that wishlist endpoints return 401 Unauthorized
// when no Authorization header is provided.
func TestWishlistRequiresAuth(t *testing.T) {
	engine := NewEngine(t)
	someID, _ := entities.NewHexID()

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/wishlist/", struct{}{}),
				netoche.WithPathParam("id", someID),
				netoche.ExpectStatus(http.StatusUnauthorized),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/wishlist/", struct{}{}),
				netoche.WithPathParam("id", someID),
				netoche.ExpectStatus(http.StatusUnauthorized),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

func TestWishlistFlowSeeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "wishlist-flow-password"

	owner := fixturesets.NewUser("wishlist-owner@example.com", userPassword, "")
	participant := fixturesets.NewUser("wishlist-participant@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(owner, participant, "Wishlist Flow Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	createReq := wishlistctrl.WishlistItemRequest{
		Label:    "Board game gift card",
		Comments: "No socks, please.",
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), participant.User.Email, userPassword),
			wishlistrunner.CreateItem(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				createReq,
				netoche.ExpectStatus(http.StatusOK),
				wishlistrunner.ExpectCreatedItem(
					wishlistctrl.WishlistItemResponse{
						Label:    createReq.Label,
						Comments: createReq.Comments,
					},
				),
			),
			wishlistrunner.ListItems(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				netoche.ExpectStatus(http.StatusOK),
				wishlistrunner.ExpectItems(
					[]wishlistctrl.WishlistItemResponse{
						{
							Label:    createReq.Label,
							Comments: createReq.Comments,
						},
					},
				),
			),
			wishlistrunner.DeleteItem(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				netoche.WithPathParamFromCtx(
					"itemId",
					func(item wishlistctrl.WishlistItemResponse) string { return item.ItemID },
				),
				netoche.ExpectStatus(http.StatusOK),
			),
			wishlistrunner.ListItems(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				netoche.ExpectStatus(http.StatusOK),
				wishlistrunner.ExpectItems([]wishlistctrl.WishlistItemResponse{}),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

// TestWishlistNonParticipantForbidden verifies that a user who is not a participant
// of the event cannot add items to the wishlist.
func TestWishlistNonParticipantForbidden(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "wishlist-forbidden-password"

	owner := fixturesets.NewUser("wishlist-forbidden-owner@example.com", userPassword, "")
	member := fixturesets.NewUser("wishlist-forbidden-member@example.com", userPassword, "")
	outsider := fixturesets.NewUser("wishlist-forbidden-outsider@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(owner, member, "Forbidden Wishlist Event")

	if err := engine.Seed(
		append(eventSet.Seedables(), outsider.User, outsider.Profile)...,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			authrunner.Login(engine.BaseURL(), outsider.User.Email, userPassword),
			wishlistrunner.CreateItem(
				engine.BaseURL(),
				eventSet.SecretFriendID(),
				wishlistctrl.WishlistItemRequest{Label: "Sneaky item"},
				netoche.ExpectStatus(http.StatusForbidden),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}

// TestWishlistCapacityExceeded verifies that adding more than 10 items to a
// participant's wishlist returns 409 Conflict.
func TestWishlistCapacityExceeded(t *testing.T) {
	engine := NewEngine(t)
	const (
		userPassword = "wishlist-cap-password"
		maxItems     = 10
	)

	owner := fixturesets.NewUser("wishlist-cap-owner@example.com", userPassword, "")
	participant := fixturesets.NewUser("wishlist-cap-participant@example.com", userPassword, "")
	eventSet := fixturesets.NewOwnerParticipant(owner, participant, "Wishlist Cap Event")

	if err := engine.Seed(eventSet.Seedables()...); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	runners := []atores.Runner{
		authrunner.Login(engine.BaseURL(), participant.User.Email, userPassword),
	}

	// Add exactly maxItems items — all must succeed.
	for range maxItems {
		runners = append(runners, wishlistrunner.CreateItem(
			engine.BaseURL(),
			eventSet.SecretFriendID(),
			wishlistctrl.WishlistItemRequest{Label: "wishlist-item"},
			netoche.ExpectStatus(http.StatusOK),
		))
	}

	// One item beyond the limit must be rejected with 409 Conflict.
	runners = append(runners, wishlistrunner.CreateItem(
		engine.BaseURL(),
		eventSet.SecretFriendID(),
		wishlistctrl.WishlistItemRequest{Label: "One too many"},
		netoche.ExpectStatus(http.StatusConflict),
	))

	mr := atores.MultiRunner{Runners: runners}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
