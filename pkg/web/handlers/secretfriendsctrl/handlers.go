package secretfriendsctrl

import (
	"context"

	"github.com/go-fuego/fuego"
	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type (
	UseCaseFactory[T any] func(ctx context.Context) (T, error)
	Controller            struct {
		web.DefaultController

		sfUseCaseFactory   UseCaseFactory[*secretfriend.UseCase]
		drawUseCaseFactory UseCaseFactory[*drawfriends.UseCase]
	}
)

func NewController(
	sfFac UseCaseFactory[*secretfriend.UseCase],
	drawFac UseCaseFactory[*drawfriends.UseCase],
) Controller {
	return Controller{
		sfUseCaseFactory:   sfFac,
		drawUseCaseFactory: drawFac,
	}
}

// CreateSecretFriend handles POST /secret-friends.
func (ctrl *Controller) CreateSecretFriend(
	c fuego.ContextWithBody[CreateSecretFriendRequest],
) (*CreateSecretFriendResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	sfUC, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	user, err := remy.GetWithContext[entities.User](
		nil,
		c.Context(),
	) // Assuming User is injected in context
	if err != nil {
		return nil, err
	}

	secretFriend, err := sfUC.Create(
		secretfriend.CreateInput{
			Name:            req.Name,
			Datetime:        req.Datetime,
			Location:        req.Location,
			MaxDenyListSize: req.MaxDenyListSize,
			OwnerID:         user.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &CreateSecretFriendResponse{
		SecretFriendID: secretFriend.ID.String(),
		InviteCode:     secretFriend.InviteCode,
	}, nil
}

// GetSecretFriend handles GET /secret-friends/{id}.
func (ctrl *Controller) GetSecretFriend(
	c fuego.ContextNoBody,
) (*GetSecretFriendResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	sfUC, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}
	secretFriend, err := sfUC.Get(id)
	if err != nil {
		return nil, err
	}

	return &GetSecretFriendResponse{
		ID:                secretFriend.ID.String(),
		Name:              secretFriend.Name,
		Datetime:          secretFriend.Datetime,
		Location:          secretFriend.Location,
		OwnerID:           secretFriend.OwnerID.String(),
		ParticipantsCount: len(secretFriend.Participants),
		Status:            string(secretFriend.Status),
	}, nil
}

// UpdateSecretFriend handles PATCH /secret-friends/{id}.
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

	sfUC, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	if err = sfUC.Update(
		secretfriend.UpdateInput{
			ID:       id,
			Name:     req.Name,
			Datetime: req.Datetime,
			Location: req.Location,
		},
	); err != nil {
		return nil, err
	}

	return map[string]any{
		"success": true,
		"message": "secret friend updated successfully",
	}, nil
}

// DrawSecretFriend handles POST /secret-friends/{id}/drawfriends.
func (ctrl *Controller) DrawSecretFriend(
	c fuego.ContextNoBody,
) (*DrawSecretFriendResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	drawUC, err := ctrl.drawUseCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	result, err := drawUC.Execute(
		drawfriends.ExecuteInput{
			SecretFriendID: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DrawSecretFriendResponse{
		SecretFriendID: id.String(),
		Status:         string(entities.StatusDrawn),
		ResultCount:    result.ParticipantCount,
	}, nil
}

// GetDrawResult handles GET /secret-friends/{id}/drawfriends-result.
func (ctrl *Controller) GetDrawResult(
	c fuego.ContextNoBody,
) (*DrawResultResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, err
	}

	drawUC, err := ctrl.drawUseCaseFactory(c.Context())
	if err != nil {
		return nil, err
	}

	user, err := remy.GetWithContext[entities.User](nil, c.Context())
	if err != nil {
		return nil, err
	}

	result, err := drawUC.GetResult(
		drawfriends.GetResultInput{
			SecretFriendID: id,
			UserID:         user.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	wishlist := make([]WishlistItem, len(result.Receiver.Wishlist.Items))
	for i, item := range result.Receiver.Wishlist.Items {
		wishlist[i] = WishlistItem{
			ItemID:   item.ID.String(),
			Label:    item.Label,
			Comments: item.Comments,
		}
	}

	return &DrawResultResponse{
		TargetUserID: result.Receiver.RelatedUser.ID.String(),
		TargetName:   result.Receiver.RelatedUser.FullName,
		Wishlist:     wishlist,
	}, nil
}
