package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type UserEditionController struct {
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
	c fuego.ContextWithBody[FormEditPassword],
) (*SuccessResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) <= 0 {
		return nil, NewHTTPError(http.StatusUnauthorized, "missing authorization token")
	}

	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// Starting with password change
	err = ctrl.serv.ChangePassword(
		token, req.CurrentPassword, req.NewPassword,
	)
	if err != nil {
		return nil, ctrl.analyseAndReturnError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}

func (ctrl UserEditionController) EditUsername(
	c fuego.ContextWithBody[FormEditUsername],
) (*SuccessResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) <= 0 {
		return nil, NewHTTPError(http.StatusUnauthorized, "missing authorization token")
	}

	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// Starting with username change
	userDTO := entities.UserBasic{
		Username: req.NewUsername,
		Password: req.CurrentPassword,
	}
	if err = ctrl.serv.ChangeUsername(token, userDTO); err != nil {
		return nil, ctrl.analyseAndReturnError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Username changed successfully",
	}, nil
}

func (ctrl UserEditionController) EditUserEmail(
	c fuego.ContextWithBody[FormEditEmail],
) (*SuccessResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if len(token) <= 0 {
		return nil, NewHTTPError(http.StatusUnauthorized, "missing authorization token")
	}

	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// Starting with email change
	userDTO := entities.UserBasic{
		Email:    req.NewEmail,
		Password: req.CurrentPassword,
	}
	if err = ctrl.serv.ChangeEmail(token, userDTO); err != nil {
		return nil, ctrl.analyseAndReturnError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "Email changed successfully",
	}, nil
}

// analyseAndReturnError analyzes the error and returns an appropriate HTTP error
func (ctrl UserEditionController) analyseAndReturnError(err error) error {
	if errors.Is(err, autherrs.ErrUserNotFound) ||
		errors.Is(err, autherrs.ErrEmailInUse) ||
		errors.Is(err, autherrs.ErrUsernameInUse) ||
		errors.Is(err, autherrs.ErrWrongPassword) {
		return NewHTTPError(http.StatusNotAcceptable, err.Error())
	}
	return err
}
