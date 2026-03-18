package wishlist

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) GetWishlist(sfID entities.HexID) ([]entities.WishlistItem, error) {
	items, err := uc.repo.GetWishlistByParticipant(
		ParticipantRef{
			UserID:         uc.associatedUser.ID,
			SecretFriendID: sfID,
		},
	)
	if err != nil {
		return nil, apperr.From(
			"wishlist_lookup_failed",
			"failed to load wishlist",
			err,
		)
	}

	return items, nil
}

func (uc *UseCase) AddItem(
	sfID entities.HexID,
	label, comments string,
) (entities.WishlistItem, error) {
	participant, err := uc.validator.CheckParticipantInSecretFriend(sfID, uc.associatedUser.ID)
	if err != nil {
		return entities.WishlistItem{}, apperr.Forbidden(
			"wishlist_access_forbidden",
			"you are not a participant in this secret friend",
			err,
		)
	}

	// Fetch current wishlist
	participantRef := ParticipantRef{
		ParticipantID:  participant.ID,
		UserID:         uc.associatedUser.ID,
		SecretFriendID: sfID,
	}

	currentList, err := uc.repo.GetWishlistByParticipant(participantRef)
	if err != nil {
		return entities.WishlistItem{}, apperr.From(
			"wishlist_lookup_failed",
			"failed to load wishlist",
			err,
		)
	}

	// Fetch sf to get config (though MaxWishListSize is not configurable by user, we can set a hardcoded limit here, e.g. 10)
	if len(currentList) >= int(uc.maxWishListSize) {
		return entities.WishlistItem{}, apperr.Conflict(
			"wishlist_capacity_reached",
			"wishlist capacity reached",
			nil,
		)
	}

	newWishItem := entities.WishlistItem{
		Label:    label,
		Comments: comments,
	}
	wishItem, err := uc.repo.AddWishlistItem(participantRef, newWishItem)
	if err != nil {
		return entities.WishlistItem{}, apperr.From(
			"wishlist_add_failed",
			"failed to add wishlist item",
			err,
		)
	}

	return wishItem, nil
}

func (uc *UseCase) DeleteItem(sfID, itemID entities.HexID) error {
	participant, err := uc.validator.CheckParticipantInSecretFriend(sfID, uc.associatedUser.ID)
	if err != nil {
		return apperr.Forbidden(
			"wishlist_access_forbidden",
			"you are not a participant in this secret friend",
			err,
		)
	}

	if err = uc.repo.RemoveWishlistItem(
		itemID,
		ParticipantRef{
			ParticipantID:  participant.ID,
			UserID:         uc.associatedUser.ID,
			SecretFriendID: sfID,
		},
	); err != nil {
		return apperr.From(
			"wishlist_remove_failed",
			"failed to remove wishlist item",
			err,
		)
	}

	return nil
}
