package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/pkg/web"
)

type UserEditionController struct {
	web.DefaultController

	serv userserv.UserEditionService
}

func NewUserEditionController(
	mailerService authserv.MailerService,
	userRepository authserv.UserAuthRepository,
	userEditionRepository userserv.UserEditionRepository,
) UserEditionController {
	return UserEditionController{
		serv: userserv.NewUserEditService(
			userRepository, userEditionRepository, mailerService,
		),
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

	err = ctrl.serv.ChangePassword.Execute(token, req.CurrentPassword, req.NewPassword)
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
	if err = ctrl.serv.ChangeUsername.Execute(token, userDTO); err != nil {
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
	if err = ctrl.serv.ChangeEmail.Execute(token, userDTO); err != nil {
		return nil, ctrl.HandleError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Email changed successfully",
	}, nil
}
