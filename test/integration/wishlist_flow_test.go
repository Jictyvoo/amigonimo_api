package integration

import (
	"net/http"
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
	wishlistrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/wishlist"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixturesets"
)

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
