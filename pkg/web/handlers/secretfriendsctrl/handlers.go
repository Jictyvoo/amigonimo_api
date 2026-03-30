package secretfriendsctrl

import (
	"context"
	"errors"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type (
	UseCaseFactory[T any] func(ctx context.Context) (T, error)
	Controller            struct {
		web.DefaultController

		sfUseCaseFactory   UseCaseFactory[secretfriend.UseCase]
		drawUseCaseFactory UseCaseFactory[drawfriends.Service]
	}
)

func NewController(
	sfFac UseCaseFactory[secretfriend.UseCase],
	drawFac UseCaseFactory[drawfriends.Service],
) Controller {
	return Controller{
		sfUseCaseFactory:   sfFac,
		drawUseCaseFactory: drawFac,
	}
}

// CreateSecretFriend handles POST /secret-friends.
func (ctrl *Controller) CreateSecretFriend(
	c fuego.Context[CreateSecretFriendRequest, struct{}],
) (*CreateSecretFriendResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	sfUC, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	secretFriend, err := sfUC.Create(
		secretfriend.CreateInput{
			Name:            req.Name,
			Datetime:        req.Datetime,
			Location:        req.Location,
			MaxDenyListSize: req.MaxDenyListSize,
		},
	)
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &CreateSecretFriendResponse{
		SecretFriendID: secretFriend.ID.String(),
		InviteCode:     secretFriend.InviteCode,
	}, nil
}

// GetSecretFriendList handles GET /.
func (ctrl *Controller) GetSecretFriendList(
	c fuego.ContextNoBody,
) (DashboardResponse, error) {
	uc, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return DashboardResponse{}, ctrl.HandleError(err)
	}

	result, err := uc.ListUserSecretFriends(entities.HexID{})
	if err != nil {
		return DashboardResponse{}, ctrl.HandleError(err)
	}

	var dashResp DashboardResponse
	parseEventList(&dashResp.Active.Created, result.Active.Created)
	parseEventList(&dashResp.Active.Participant, result.Active.Participant)
	parseEventList(&dashResp.Inactive.Created, result.Inactive.Created)
	parseEventList(&dashResp.Inactive.Participant, result.Inactive.Participant)

	return dashResp, nil
}

// GetSecretFriend handles GET /secret-friends/{id}.
func (ctrl *Controller) GetSecretFriend(
	c fuego.ContextNoBody,
) (*GetSecretFriendResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	sfUC, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}
	secretFriend, err := sfUC.Get(id)
	if err != nil {
		return nil, ctrl.HandleError(err)
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
	c fuego.Context[UpdateSecretFriendRequest, struct{}],
) (UpdateSecretFriendResponse, error) {
	req, err := c.Body()
	if err != nil {
		return UpdateSecretFriendResponse{}, ctrl.HandleError(err)
	}

	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return UpdateSecretFriendResponse{}, ctrl.HandleError(err)
	}

	sfUC, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return UpdateSecretFriendResponse{}, ctrl.HandleError(err)
	}

	if err = sfUC.Update(
		secretfriend.UpdateInput{
			ID:       id,
			Name:     req.Name,
			Datetime: req.Datetime,
			Location: req.Location,
		},
	); err != nil {
		return UpdateSecretFriendResponse{}, ctrl.HandleError(err)
	}

	return UpdateSecretFriendResponse{
		Success: true,
		Message: "secret friend updated successfully",
	}, nil
}

// DrawSecretFriend handles POST /secret-friends/{id}/draw.
func (ctrl *Controller) DrawSecretFriend(
	c fuego.ContextNoBody,
) (*DrawSecretFriendResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	drawUC, err := ctrl.drawUseCaseFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	result, err := drawUC.Execute(id)
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &DrawSecretFriendResponse{
		SecretFriendID: id.String(),
		Status:         string(entities.StatusDrawn),
		ResultCount:    result,
	}, nil
}

// GetDrawResult handles GET /secret-friends/{id}/drawfriends-result.
func (ctrl *Controller) GetDrawResult(
	c fuego.ContextNoBody,
) (*DrawResultResponse, error) {
	id, err := ctrl.ParamID(c.Request())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	drawUC, err := ctrl.drawUseCaseFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	result, err := drawUC.GetResult(id)
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	wishlist := make([]WishlistItem, len(result.Receiver.Wishlist))
	for i, item := range result.Receiver.Wishlist {
		wishlist[i] = WishlistItem{
			ItemID:   item.ID.String(),
			Label:    item.Label,
			Comments: item.Comments,
		}
	}

	return &DrawResultResponse{
		TargetUserID: result.Receiver.RelatedUser.ID.String(),
		TargetName:   result.Receiver.Profile.FullName,
		Wishlist:     wishlist,
	}, nil
}

// GetInviteByCode handles GET /invites/{code}.
func (ctrl *Controller) GetInviteByCode(
	c fuego.ContextNoBody,
) (*InviteInfoResponse, error) {
	code := c.PathParam("code")
	if code == "" {
		return nil, ctrl.HTTPError(400, errors.New("invite code is required"))
	}

	uc, err := ctrl.sfUseCaseFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	sf, err := uc.GetInviteInfo(code)
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &InviteInfoResponse{
		SecretFriendID: sf.ID.String(),
		Name:           sf.Name,
	}, nil
}
