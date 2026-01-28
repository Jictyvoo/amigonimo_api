package denylistctrl

import "github.com/jictyvoo/amigonimo_api/internal/entities"

func parseDeniedUser(item entities.DeniedUser) DeniedUserResponse {
	return DeniedUserResponse{
		UserID:   item.InnerParticipant.RelatedUser.ID.String(),
		Fullname: item.InnerParticipant.RelatedUser.FullName,
	}
}
