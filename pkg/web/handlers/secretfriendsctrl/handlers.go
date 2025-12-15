package secretfriendsctrl

import (
	"context"
	"net/http"

	"github.com/go-fuego/fuego"
	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services/evtserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type (
	ServiceFactory func(ctx context.Context) (*evtserv.Service, error)
	Controller     struct {
		serviceFactory ServiceFactory
		web.DefaultController
	}
)

func NewController(servFac ServiceFactory) Controller {
	return Controller{serviceFactory: servFac}
}

// CreateSecretFriend handles POST /secret-friends
func (ctrl *Controller) CreateSecretFriend(
	c fuego.ContextWithBody[CreateSecretFriendRequest],
) (*CreateSecretFriendResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	serv, err := ctrl.serviceFactory(c.Context())
	if err != nil {
		return nil, err
	}
	secretFriend, err := serv.CreateSecretFriend(
		req.Name,
		req.Datetime,
		req.Location,
		req.MaxDenyListSize,
	)
	if err != nil {
		return nil, err
	}

	return &CreateSecretFriendResponse{
		SecretFriendID: secretFriend.ID.String(),
		InviteCode:     secretFriend.InviteCode,
		InviteLink:     secretFriend.InviteLink,
	}, nil
}

// GetSecretFriend handles GET /secret-friends/{id}
func (ctrl *Controller) GetSecretFriend(
	c fuego.ContextNoBody,
) (*GetSecretFriendResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	serv, err := ctrl.serviceFactory(c.Context())
	if err != nil {
		return nil, err
	}
	secretFriend, err := serv.GetSecretFriend(id)
	if err != nil {
		return nil, err
	}

	// Get participants count (assuming service provides this)
	// For now, we'll use the length of participants array
	participantsCount := len(secretFriend.Participants)

	return &GetSecretFriendResponse{
		ID:                secretFriend.ID.String(),
		Name:              secretFriend.Name,
		Datetime:          secretFriend.Datetime,
		Location:          secretFriend.Location,
		OwnerID:           secretFriend.OwnerID.String(),
		ParticipantsCount: participantsCount,
		Status:            string(secretFriend.Status),
	}, nil
}

// UpdateSecretFriend handles PATCH /secret-friends/{id}
func (ctrl *Controller) UpdateSecretFriend(
	c fuego.ContextWithBody[UpdateSecretFriendRequest],
) (any, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	serv, err := ctrl.serviceFactory(c.Context())
	if err != nil {
		return nil, err
	}
	if err = serv.UpdateSecretFriend(
		id,
		req.Name,
		req.Datetime,
		req.Location,
	); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "secret friend updated successfully",
	}, nil
}

// DrawSecretFriend handles POST /secret-friends/{id}/draw
func (ctrl *Controller) DrawSecretFriend(
	c fuego.ContextNoBody,
) (*DrawSecretFriendResponse, error) {
	idStr := c.Request().PathValue("id")
	if idStr == "" {
		return nil, ctrl.HTTPError(http.StatusBadRequest, "id is required")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, ctrl.HTTPError(http.StatusBadRequest, "invalid id format")
	}

	serv, err := ctrl.serviceFactory(c.Context())
	if err != nil {
		return nil, err
	}
	resultCount, err := serv.DrawSecretFriend(entities.HexID(id))
	if err != nil {
		return nil, err
	}

	return &DrawSecretFriendResponse{
		SecretFriendID: idStr,
		Status:         string(entities.StatusDrawn),
		ResultCount:    resultCount,
	}, nil
}

// GetDrawResult handles GET /secret-friends/{id}/draw-result
func (ctrl *Controller) GetDrawResult(
	c fuego.ContextNoBody,
) (*DrawResultResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, err
	}
	serv, err := ctrl.serviceFactory(c.Context())
	if err != nil {
		return nil, err
	}
	result, err := serv.GetDrawResultForUser(id)
	if err != nil {
		return nil, err
	}

	// Map wishlist items from participant
	wishlist := make([]WishlistItem, 0)
	// Note: Wishlist items would need to be fetched from the Receiver participant
	// This is a placeholder for the actual wishlist mapping

	return &DrawResultResponse{
		TargetUserID: result.Receiver.RelatedUser.ID.String(),
		TargetName:   result.Receiver.RelatedUser.FullName,
		Wishlist:     wishlist,
	}, nil
}
