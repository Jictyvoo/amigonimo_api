package wishlist

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Repository interface {
	AddWishlistItem(
		participant entities.Participant,
		wishItem entities.WishlistItem,
	) (entities.WishlistItem, error)
	RemoveWishlistItem(itemID entities.HexID, participant entities.Participant) error
	GetWishlistByParticipant(participant entities.Participant) ([]entities.WishlistItem, error)
}

type UseCase struct {
	repo           Repository
	associatedUser entities.User
}

func New(associatedUser entities.User, repo Repository) UseCase {
	return UseCase{associatedUser: associatedUser, repo: repo}
}
