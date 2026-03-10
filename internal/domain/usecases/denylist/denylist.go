package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Repository interface {
	AddDenyListEntry(
		p entities.Participant,
		deniedUserID entities.HexID,
	) (entities.DeniedUser, error)
	RemoveDenyListEntry(p entities.Participant, deniedUserID entities.HexID) error
	GetDenyListByParticipant(p entities.Participant) ([]entities.DeniedUser, error)
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
