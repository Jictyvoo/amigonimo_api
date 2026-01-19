package secretfriendsctrl

import "github.com/jictyvoo/amigonimo_api/internal/entities"

func parseEventList(dest *[]SecretFriendSummary, source []entities.SecretFriend) {
	*dest = make([]SecretFriendSummary, len(source))
	for i, sf := range source {
		(*dest)[i] = mapToSummary(sf)
	}
}

func mapToSummary(sf entities.SecretFriend) SecretFriendSummary {
	return SecretFriendSummary{
		ID:                sf.ID.String(),
		Name:              sf.Name,
		Datetime:          sf.Datetime,
		Location:          sf.Location,
		Status:            string(sf.Status),
		ParticipantsCount: uint8(len(sf.Participants)),
	}
}
