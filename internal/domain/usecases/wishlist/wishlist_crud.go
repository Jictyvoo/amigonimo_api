package wishlist

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) GetWishlist(sfID entities.HexID) ([]entities.WishlistItem, error) {
	participant := entities.NewParticipant(sfID, uc.associatedUser)
	return uc.repo.GetWishlistByParticipant(participant)
}

func (uc *UseCase) AddItem(
	sfID entities.HexID,
	label, comments string,
) (entities.WishlistItem, error) {
	participant := entities.NewParticipant(sfID, uc.associatedUser)
	newWishItem := entities.WishlistItem{
		Label:    label,
		Comments: comments,
	}
	return uc.repo.AddWishlistItem(participant, newWishItem)
}

func (uc *UseCase) DeleteItem(sfID, itemID entities.HexID) error {
	participant := entities.NewParticipant(sfID, uc.associatedUser)
	return uc.repo.RemoveWishlistItem(itemID, participant)
}
