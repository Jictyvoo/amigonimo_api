package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Repository interface {
	AddDenyListEntry(
		participant ParticipantRef,
		deniedUserID entities.HexID,
	) (entities.DeniedUser, error)
	RemoveDenyListEntry(participant ParticipantRef, deniedUserID entities.HexID) error
	GetDenyListByParticipant(participant ParticipantRef) ([]entities.DeniedUser, error)
}

type ParticipantRef struct {
	ParticipantID  entities.HexID
	UserID         entities.HexID
	SecretFriendID entities.HexID
}

type UseCase struct {
	repo           Repository
	facProvider    FacadeProvider
	associatedUser entities.User
}

func New(
	associatedUser entities.User,
	repo Repository,
	provider FacadeProvider,
) UseCase {
	return UseCase{
		associatedUser: associatedUser,
		repo:           repo,
		facProvider:    provider,
	}
}
