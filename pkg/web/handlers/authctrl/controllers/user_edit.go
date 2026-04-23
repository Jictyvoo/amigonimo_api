package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type UserEditionServiceFactory func(ctx context.Context) (userserv.UserEditionService, error)

type UserEditionController struct {
	web.DefaultController

	servFactory UserEditionServiceFactory
}

func NewUserEditionController(
	servFactory UserEditionServiceFactory,
) UserEditionController {
	return UserEditionController{
		servFactory: servFactory,
	}
}

func (ctrl UserEditionController) EditUserPassword(
	c fuego.Context[FormEditPassword, struct{}],
) (*SuccessResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) <= 0 {
		return nil, ctrl.HTTPError(
			http.StatusUnauthorized,
			errors.New("missing authorization token"),
		)
	}

	req, err := c.Body()
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	serv, err := ctrl.servFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	err = serv.ChangePassword.Execute(token, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}

func (ctrl UserEditionController) EditUsername(
	c fuego.Context[FormEditUsername, struct{}],
) (*SuccessResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) <= 0 {
		return nil, ctrl.HTTPError(
			http.StatusUnauthorized,
			errors.New("missing authorization token"),
		)
	}

	req, err := c.Body()
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	userDTO := authvalues.UserBasic{
		Username: req.NewUsername,
		Password: req.CurrentPassword,
	}
	serv, err := ctrl.servFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	if err = serv.ChangeUsername.Execute(token, userDTO); err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Username changed successfully",
	}, nil
}

func (ctrl UserEditionController) EditUserEmail(
	c fuego.Context[FormEditEmail, struct{}],
) (*SuccessResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) <= 0 {
		return nil, ctrl.HTTPError(
			http.StatusUnauthorized,
			errors.New("missing authorization token"),
		)
	}

	req, err := c.Body()
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	userDTO := authvalues.UserBasic{
		Email:    req.NewEmail,
		Password: req.CurrentPassword,
	}
	serv, err := ctrl.servFactory(c.Context())
	if err != nil {
		return nil, ctrl.HandleError(err)
	}

	if err = serv.ChangeEmail.Execute(token, userDTO); err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Email changed successfully",
	}, nil
}
