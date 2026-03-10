package wishlist

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// ParticipantFacade defines what wishlist module needs from participant.
type ParticipantFacade interface {
	ports.Facade
	CheckParticipantInSecretFriend(sfID, userID entities.HexID) (entities.Participant, error)
}

type Repository interface {
	AddWishlistItem(
		participant entities.Participant,
		wishItem entities.WishlistItem,
	) (entities.WishlistItem, error)
	RemoveWishlistItem(itemID entities.HexID, participant entities.Participant) error
	GetWishlistByParticipant(participant entities.Participant) ([]entities.WishlistItem, error)
}

type UseCase struct {
	repo            Repository
	validator       ParticipantFacade
	associatedUser  entities.User
	maxWishListSize uint8
}

func New(
	associatedUser entities.User,
	repo Repository,
	participantFacade ParticipantFacade,
) UseCase {
	return UseCase{
		associatedUser:  associatedUser,
		repo:            repo,
		validator:       participantFacade,
		maxWishListSize: 10,
	}
}
