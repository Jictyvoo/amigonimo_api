package integration

import (
	"net/http"
	"testing"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

func TestWishlistFlowSeeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "wishlist-flow-password"

	ownerBuilder := fixtures.NewUser().
		WithEmail("wishlist-owner@example.com").
		WithPassword(userPassword)
	owner := ownerBuilder.Build()
	ownerProfile := ownerBuilder.BuildProfile()
	participantBuilder := fixtures.NewUser().
		WithEmail("wishlist-participant@example.com").
		WithPassword(userPassword)
	participant := participantBuilder.Build()
	participantProfile := participantBuilder.BuildProfile()
	secretFriend := fixtures.NewSecretFriend().
		WithOwner(owner).
		WithName("Wishlist Flow Event").
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
		participantEntry,
	); err != nil {
		t.Fatalf("seedErr: %v", err)
	}

	secretFriendID, _ := entities.NewHexIDFromBytes(secretFriend.ID)
	createReq := wishlistctrl.WishlistItemRequest{
		Label:    "Board game gift card",
		Comments: "No socks, please.",
	}

	mr := atores.MultiRunner{
		Runners: []atores.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), participant.Email, userPassword),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/wishlist/", createReq),
				netoche.WithPathParam("id", secretFriendID),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					wishlistctrl.WishlistItemResponse{
						Label:    createReq.Label,
						Comments: createReq.Comments,
					},
					func(expected, actual *wishlistctrl.WishlistItemResponse) error {
						expected.ItemID = actual.ItemID
						expected.AddedAt = actual.AddedAt
						return nil
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/wishlist/", struct{}{}),
				netoche.WithPathParam("id", secretFriendID),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					[]wishlistctrl.WishlistItemResponse{
						{
							Label:    createReq.Label,
							Comments: createReq.Comments,
						},
					},
					func(expected, actual *[]wishlistctrl.WishlistItemResponse) error {
						limit := min(len(*actual), len(*expected))
						for index := 0; index < limit; index++ {
							(*expected)[index].ItemID = (*actual)[index].ItemID
							(*expected)[index].AddedAt = (*actual)[index].AddedAt
						}
						return nil
					},
				),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(
					http.MethodDelete, "/secret-friends/{id}/wishlist/{itemId}", struct{}{},
				),
				netoche.WithPathParam("id", secretFriendID),
				netoche.WithPathParamFromCtx(
					"itemId",
					func(item wishlistctrl.WishlistItemResponse) string { return item.ItemID },
				),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusOK),
			),
			netoche.New(
				engine.BaseURL(),
				netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/wishlist/", struct{}{}),
				netoche.WithPathParam("id", secretFriendID),
				stdrunners.WithAuthHeaderFromLogin(),
				netoche.ExpectStatus(http.StatusOK),
				netoche.ExpectBody(
					[]wishlistctrl.WishlistItemResponse{},
				),
			),
		},
	}

	if err := engine.Execute(t, mr); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
