package invitesctrl

import (
	"github.com/go-fuego/fuego"
)

type InvitesHandlers struct {
	// TODO: Add service dependencies
}

func NewInvitesHandlers() *InvitesHandlers {
	return &InvitesHandlers{}
}

// GetInviteByCode handles GET /invites/{code}.
func (h *InvitesHandlers) GetInviteByCode(
	c fuego.ContextNoBody,
) (*InviteInfoResponse, error) {
	// TODO: Extract code from path
	// TODO: Implement service call
	return nil, nil
}
