package secretfriendsctrl

import (
	"github.com/go-fuego/fuego"
)

type SecretFriendsHandlers struct {
	// TODO: Add service dependencies
}

func NewSecretFriendsHandlers() *SecretFriendsHandlers {
	return &SecretFriendsHandlers{}
}

// CreateSecretFriend handles POST /secret-friends.
func (h *SecretFriendsHandlers) CreateSecretFriend(
	c fuego.ContextWithBody[CreateSecretFriendRequest],
) (*CreateSecretFriendResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// TODO: Implement service call
	_ = req
	return nil, nil
}

// GetSecretFriend handles GET /secret-friends/{id}.
func (h *SecretFriendsHandlers) GetSecretFriend(
	c fuego.ContextNoBody,
) (*GetSecretFriendResponse, error) {
	// TODO: Extract id from path
	// TODO: Implement service call
	return nil, nil
}

// UpdateSecretFriend handles PATCH /secret-friends/{id}.
func (h *SecretFriendsHandlers) UpdateSecretFriend(
	c fuego.ContextWithBody[UpdateSecretFriendRequest],
) (any, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// TODO: Extract id from path
	// TODO: Implement service call
	_ = req
	return nil, nil
}

// DrawSecretFriend handles POST /secret-friends/{id}/draw.
func (h *SecretFriendsHandlers) DrawSecretFriend(
	c fuego.ContextNoBody,
) (*DrawSecretFriendResponse, error) {
	// TODO: Extract id from path
	// TODO: Implement service call
	return nil, nil
}

// GetDrawResult handles GET /secret-friends/{id}/draw-result.
func (h *SecretFriendsHandlers) GetDrawResult(
	c fuego.ContextNoBody,
) (*DrawResultResponse, error) {
	// TODO: Extract secretFriendId from path
	// TODO: Implement service call
	return nil, nil
}
