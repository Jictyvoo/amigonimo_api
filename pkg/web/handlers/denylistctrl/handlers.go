package denylistctrl

import (
	"github.com/go-fuego/fuego"
)

type DenyListHandlers struct {
	// TODO: Add service dependencies
}

func NewDenyListHandlers() *DenyListHandlers {
	return &DenyListHandlers{}
}

// GetDenyList handles GET /secret-friends/{id}/denylist
func (h *DenyListHandlers) GetDenyList(
	c fuego.ContextNoBody,
) ([]DeniedUserResponse, error) {
	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	return nil, nil
}

// AddDenyListEntry handles POST /secret-friends/{id}/denylist
func (h *DenyListHandlers) AddDenyListEntry(
	c fuego.ContextWithBody[AddDenyListRequest],
) (DeniedUserResponse, error) {
	req, err := c.Body()
	if err != nil {
		return DeniedUserResponse{}, err
	}

	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	_ = req
	return DeniedUserResponse{}, nil
}

// RemoveDenyListEntry handles DELETE /secret-friends/{id}/denylist/{targetUserId}
func (h *DenyListHandlers) RemoveDenyListEntry(
	c fuego.ContextNoBody,
) (any, error) {
	// TODO: Extract secretFriendId and targetUserId from path
	// TODO: Implement service call
	return nil, nil
}
