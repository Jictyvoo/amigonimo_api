package wishlistrunner

import (
	"net/http"

	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores"
	"github.com/wrapped-owls/testereiro/puppetest/pkg/atores/netoche"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web/handlers/wishlistctrl"
	authrunner "github.com/jictyvoo/amigonimo_api/test/integration/runners/auth"
)

func CreateItem(
	baseURL string,
	secretFriendID entities.HexID,
	req wishlistctrl.WishlistItemRequest,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodPost, "/secret-friends/{id}/wishlist/", req),
		netoche.WithPathParam("id", secretFriendID),
		authrunner.WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func ListItems(
	baseURL string,
	secretFriendID entities.HexID,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(http.MethodGet, "/secret-friends/{id}/wishlist/", struct{}{}),
		netoche.WithPathParam("id", secretFriendID),
		authrunner.WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func DeleteItem(
	baseURL string,
	secretFriendID entities.HexID,
	opts ...netoche.Option,
) atores.Runner {
	baseOpts := []netoche.Option{
		netoche.WithRequest(
			http.MethodDelete,
			"/secret-friends/{id}/wishlist/{itemId}",
			struct{}{},
		),
		netoche.WithPathParam("id", secretFriendID),
		authrunner.WithAuthHeaderFromLogin(),
	}

	return netoche.New(baseURL, append(baseOpts, opts...)...)
}

func ExpectCreatedItem(expected wishlistctrl.WishlistItemResponse) netoche.Option {
	return netoche.ExpectBody(
		expected,
		func(exp, actual *wishlistctrl.WishlistItemResponse) error {
			exp.ItemID = actual.ItemID
			exp.AddedAt = actual.AddedAt
			return nil
		},
	)
}

func ExpectItems(expected []wishlistctrl.WishlistItemResponse) netoche.Option {
	return netoche.ExpectBody(
		expected,
		func(exp, actual *[]wishlistctrl.WishlistItemResponse) error {
			limit := min(len(*actual), len(*exp))
			for i := 0; i < limit; i++ {
				(*exp)[i].ItemID = (*actual)[i].ItemID
				(*exp)[i].AddedAt = (*actual)[i].AddedAt
			}
			return nil
		},
	)
}
