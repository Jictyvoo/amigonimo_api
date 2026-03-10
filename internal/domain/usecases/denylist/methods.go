package denylist

import (
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc UseCase) GetDenyList(sfID entities.HexID) ([]entities.DeniedUser, error) {
	return uc.repo.GetDenyListByParticipant(
		entities.NewParticipant(sfID, uc.associatedUser),
	)
}

func (uc UseCase) AddEntry(sfID, deniedUserID entities.HexID) (entities.DeniedUser, error) {
	if uc.associatedUser.ID == deniedUserID {
		return entities.DeniedUser{}, fmt.Errorf("cannot add yourself to the denylist")
	}

	participant, err := uc.facProvider.participant.CheckParticipantInSecretFriend(
		sfID, uc.associatedUser.ID,
	)
	if err != nil {
		return entities.DeniedUser{}, fmt.Errorf("user is not a participant: %w", err)
	}

	if _, err = uc.facProvider.participant.CheckParticipantInSecretFriend(sfID, deniedUserID); err != nil {
		return entities.DeniedUser{}, fmt.Errorf("target user is not a participant: %w", err)
	}

	// Validate capacity
	currentList, err := uc.repo.GetDenyListByParticipant(participant)
	if err != nil {
		return entities.DeniedUser{}, fmt.Errorf("failed to get current denylist: %w", err)
	}

	sf, err := uc.facProvider.secretFriend.GetSecretFriendByID(sfID)
	if err != nil {
		return entities.DeniedUser{}, fmt.Errorf("failed to fetch secret-friend config: %w", err)
	}
	if len(currentList) >= int(sf.MaxDenyListSize) {
		return entities.DeniedUser{}, fmt.Errorf(
			"denylist capacity reached: max %d",
			sf.MaxDenyListSize,
		)
	}

	return uc.repo.AddDenyListEntry(participant, deniedUserID)
}

func (uc UseCase) RemoveEntry(sfID, deniedUserID entities.HexID) error {
	return uc.repo.RemoveDenyListEntry(
		entities.NewParticipant(sfID, uc.associatedUser), deniedUserID,
	)
}
