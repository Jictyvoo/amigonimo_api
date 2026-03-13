package wishlist

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) GetWishlist(sfID entities.HexID) ([]entities.WishlistItem, error) {
	return uc.repo.GetWishlistByParticipant(
		ParticipantRef{
			UserID:         uc.associatedUser.ID,
			SecretFriendID: sfID,
		},
	)
}

func (uc *UseCase) AddItem(
	sfID entities.HexID,
	label, comments string,
) (entities.WishlistItem, error) {
	participant, err := uc.validator.CheckParticipantInSecretFriend(sfID, uc.associatedUser.ID)
	if err != nil {
		return entities.WishlistItem{}, fmt.Errorf("validation failed: %w", err)
	}

	// Fetch current wishlist
	participantRef := ParticipantRef{
		ParticipantID:  participant.ID,
		UserID:         uc.associatedUser.ID,
		SecretFriendID: sfID,
	}

	currentList, err := uc.repo.GetWishlistByParticipant(participantRef)
	if err != nil {
		return entities.WishlistItem{}, fmt.Errorf("could not get current wishlist: %w", err)
	}

	// Fetch sf to get config (though MaxWishListSize is not configurable by user, we can set a hardcoded limit here, e.g. 10)
	if len(currentList) >= int(uc.maxWishListSize) {
		return entities.WishlistItem{}, fmt.Errorf(
			"wishlist capacity reached: max %d", uc.maxWishListSize,
		)
	}

	newWishItem := entities.WishlistItem{
		Label:    label,
		Comments: comments,
	}
	return uc.repo.AddWishlistItem(participantRef, newWishItem)
}

func (uc *UseCase) DeleteItem(sfID, itemID entities.HexID) error {
	participant, err := uc.validator.CheckParticipantInSecretFriend(sfID, uc.associatedUser.ID)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return uc.repo.RemoveWishlistItem(
		itemID,
		ParticipantRef{
			ParticipantID:  participant.ID,
			UserID:         uc.associatedUser.ID,
			SecretFriendID: sfID,
		},
	)
}
