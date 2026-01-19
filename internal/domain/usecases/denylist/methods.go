package denylist

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc UseCase) GetDenyList(sfID entities.HexID) ([]entities.DeniedUser, error) {
	return uc.repo.GetDenyListByParticipant(
		entities.NewParticipant(sfID, uc.associatedUser),
	)
}

func (uc UseCase) AddEntry(sfID, deniedUserID entities.HexID) (entities.DeniedUser, error) {
	return uc.repo.AddDenyListEntry(
		entities.NewParticipant(sfID, uc.associatedUser), deniedUserID,
	)
}

func (uc UseCase) RemoveEntry(sfID, deniedUserID entities.HexID) error {
	return uc.repo.RemoveDenyListEntry(
		entities.NewParticipant(sfID, uc.associatedUser), deniedUserID,
	)
}
