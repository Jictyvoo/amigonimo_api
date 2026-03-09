package participant

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Repository interface {
	AddParticipant(sfID, userID entities.HexID) (entities.Participant, error)
	ListParticipants(sfID entities.HexID) ([]entities.Participant, error)
	RemoveParticipant(sfID, userID entities.HexID) error
}

type UseCase struct {
	repo           Repository
	associatedUser entities.User
}

func New(associatedUser entities.User, repo Repository) UseCase {
	return UseCase{associatedUser: associatedUser, repo: repo}
}
