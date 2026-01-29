package integration

import (
	"net/http"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
	"github.com/jictyvoo/amigonimo_api/test/integration/stdrunners"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners"
	"github.com/jictyvoo/amigonimo_api/test/internal/runners/reqrunner"
)

func TestWishlistFlowSeeded(t *testing.T) {
	engine := NewEngine(t)
	const userPassword = "wishlist-flow-password"

	owner := fixtures.NewUser().
		WithEmail("wishlist-owner@example.com").
		WithPassword(userPassword).
		Build()
	participant := fixtures.NewUser().
		WithEmail("wishlist-participant@example.com").
		WithPassword(userPassword).
		Build()
	secretFriend := fixtures.NewSecretFriend().
		WithOwner(owner).
		WithName("Wishlist Flow Event").
		Build()
	participantEntry := fixtures.NewParticipant().
		WithUser(participant).
		WithSecretFriend(secretFriend).
		Build()

	engine.Seed(owner, participant, secretFriend, participantEntry)

	secretFriendID, _ := entities.NewHexIDFromBytes(secretFriend.ID)
	createReq := wishlistctrl.WishlistItemRequest{
		Label:    "Board game gift card",
		Comments: "No socks, please.",
	}

	mr := runners.MultiRunner{
		Runners: []runners.Runner{
			stdrunners.LoginRunner(engine.BaseURL(), participant.Email, userPassword),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(http.MethodPost, "/secret-friends/{id}/wishlist/", createReq),
				reqrunner.WithPathParam("id", secretFriendID),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
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
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(http.MethodGet, "/secret-friends/{id}/wishlist/", struct{}{}),
				reqrunner.WithPathParam("id", secretFriendID),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					[]wishlistctrl.WishlistItemResponse{
						{
							Label:    createReq.Label,
							Comments: createReq.Comments,
						},
					},
					func(expected, actual *[]wishlistctrl.WishlistItemResponse) error {
						limit := len(*expected)
						if len(*actual) < limit {
							limit = len(*actual)
						}
						for index := 0; index < limit; index++ {
							(*expected)[index].ItemID = (*actual)[index].ItemID
							(*expected)[index].AddedAt = (*actual)[index].AddedAt
						}
						return nil
					},
				),
			),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(
					http.MethodDelete, "/secret-friends/{id}/wishlist/{itemId}", struct{}{},
				),
				reqrunner.WithPathParam("id", secretFriendID),
				reqrunner.WithPathParamFromCtx(
					"itemId",
					func(item wishlistctrl.WishlistItemResponse) string { return item.ItemID },
				),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.ExpectStatus(http.StatusOK),
			),
			reqrunner.NewHttpRunner(
				engine.BaseURL(),
				reqrunner.WithRequest(http.MethodGet, "/secret-friends/{id}/wishlist/", struct{}{}),
				reqrunner.WithPathParam("id", secretFriendID),
				stdrunners.WithAuthHeaderFromLogin(),
				reqrunner.ExpectStatus(http.StatusOK),
				reqrunner.ExpectBody(
					[]wishlistctrl.WishlistItemResponse{},
				),
			),
		},
	}

	if err := mr.Run(t); err != nil {
		t.Fatalf("MultiRunner failed: %v", err)
	}
}
