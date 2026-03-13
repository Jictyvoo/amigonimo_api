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
		participant ParticipantRef,
		wishItem entities.WishlistItem,
	) (entities.WishlistItem, error)
	RemoveWishlistItem(itemID entities.HexID, participant ParticipantRef) error
	GetWishlistByParticipant(participant ParticipantRef) ([]entities.WishlistItem, error)
}

type ParticipantRef struct {
	ParticipantID  entities.HexID
	UserID         entities.HexID
	SecretFriendID entities.HexID
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
