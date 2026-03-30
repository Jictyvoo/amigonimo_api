package denylistctrl

import "github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"

func parseDeniedUser(item denylist.DeniedEntry) DeniedUserResponse {
	return DeniedUserResponse{
		UserID:   item.DeniedUserID.String(),
		Fullname: item.FullName,
	}
}
