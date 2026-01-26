package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (h *AuthenticationController) SignUp(
	c fuego.ContextWithBody[FormUser],
) (*SuccessResponse, error) {
	req, err := c.Body()
	if err != nil {
		return &SuccessResponse{Message: "Failed to obtain request body"}, err
	}

	userDTO := entities.UserBasic{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if _, err = h.authServ.UserSignUp(userDTO); err != nil {
		if errors.Is(err, autherrs.ErrEmailUsed) {
			return nil, NewHTTPError(http.StatusPreconditionFailed, err.Error())
		}
		return nil, err
	}

	c.SetStatus(http.StatusCreated)
	return &SuccessResponse{Success: true, Message: "User created successfully"}, nil
}

func (h *AuthenticationController) LogIn(
	c fuego.ContextWithBody[FormUser],
) (*LoginResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	userDTO := entities.UserBasic{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	authToken, err := h.authServ.UserLogIn(userDTO)
	if err != nil {
		if errors.Is(err, autherrs.ErrUserEmailNotFound) ||
			errors.Is(err, autherrs.ErrWrongPassword) {
			return nil, NewHTTPError(
				http.StatusNotAcceptable,
				autherrs.ErrUserEmailNotFound.Error(),
			)
		}
		return nil, err
	}

	// Generate JWT from AuthenticationToken
	jwtToken, err := h.generateJWT(authToken)
	if err != nil {
		return nil, err
	}

	// Set Authorization header with token
	c.Response().Header().Set("Authorization", jwtToken)

	return &LoginResponse{Token: jwtToken}, nil
}

func (h *AuthenticationController) RegenerateAuthToken(
	c fuego.ContextNoBody,
) (*LoginResponse, error) {
	authHeader := c.Request().Header.Get("Authorization")
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

	authToken, err := h.authServ.RegenerateLogin(refreshToken)
	if err != nil {
		if errors.Is(err, autherrs.ErrInvalidAuthToken) {
			return nil, NewHTTPError(http.StatusPreconditionFailed, err.Error())
		}
		return nil, err
	}

	// Generate JWT from AuthenticationToken
	jwtToken, err := h.generateJWT(authToken)
	if err != nil {
		return nil, err
	}

	// Set Authorization header with new token
	c.Response().Header().Set("Authorization", jwtToken)

	return &LoginResponse{Token: jwtToken}, nil
}

func (h *AuthenticationController) VerifyUser(
	c fuego.ContextNoBody,
) (*SuccessResponse, error) {
	verifyCode := c.Request().PathValue("verify_code")
	if verifyCode == "" {
		return nil, NewHTTPError(http.StatusPreconditionRequired, "verification code is required")
	}

	decoded, err := base64.StdEncoding.DecodeString(verifyCode)
	if err != nil {
		return nil, NewHTTPError(
			http.StatusPreconditionRequired,
			"invalid verification code format",
		)
	}

	if err = h.authServ.VerifyUserCode(string(decoded)); err != nil {
		return nil, NewHTTPError(http.StatusPreconditionFailed, err.Error())
	}

	return &SuccessResponse{
		Success: true,
		Message: "User successfully verified",
	}, nil
}

func (h *AuthenticationController) ForgotPassword(
	c fuego.ContextWithBody[FormUser],
) (*ForgotPasswordResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	if len(req.Username) > 0 {
		obfuscated := h.authServ.GetObfuscatedEmail(req.Username)
		if len(obfuscated) > 0 {
			c.Response().WriteHeader(http.StatusNonAuthoritativeInfo)
			return &ForgotPasswordResponse{
				SuccessResponse: SuccessResponse{Success: true},
				ObfuscatedEmail: obfuscated,
			}, nil
		}
		return nil, NewHTTPError(http.StatusNotAcceptable, "user not found")
	}

	if len(req.Email) > 0 {
		if err = h.authServ.GeneratePasswordRecovery(req.Email); err != nil {
			if errors.Is(err, autherrs.ErrUserEmailNotFound) {
				return nil, NewHTTPError(http.StatusPreconditionFailed, err.Error())
			}
			return nil, err
		}
		// The client will open a window to put this code, and send a request to server again with new password
		c.Response().WriteHeader(http.StatusAccepted)
		return &ForgotPasswordResponse{
			SuccessResponse: SuccessResponse{
				Success: true,
				Message: "an email with a recover code was sent",
			},
		}, nil
	}

	return nil, NewHTTPError(http.StatusUnprocessableEntity, "no user or email provided")
}

func (h *AuthenticationController) CheckRecoveryCode(
	c fuego.ContextWithBody[FormRecoveryCode],
) (*SuccessResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	if _, err = h.authServ.CheckRecoveryCode(req.Email, req.RecoveryCode); err != nil {
		if errors.Is(err, autherrs.ErrUserRecoveryNotFound) {
			return nil, NewHTTPError(http.StatusPreconditionFailed, err.Error())
		}
		return nil, err
	}

	return &SuccessResponse{
		Success: true,
		Message: "sent recovery-code is valid",
	}, nil
}

func (h *AuthenticationController) PasswordReset(
	c fuego.ContextWithBody[FormResetPassword],
) (*SuccessResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	userDTO := entities.UserBasic{
		Email:    req.Email,
		Password: req.NewPassword,
	}

	if err = h.authServ.ResetPassword(userDTO, req.RecoveryCode); err != nil {
		if errors.Is(err, autherrs.ErrUserRecoveryNotFound) {
			return nil, NewHTTPError(http.StatusPreconditionFailed, err.Error())
		}
		return nil, err
	}

	c.Response().WriteHeader(http.StatusCreated)
	return &SuccessResponse{
		Success: true,
		Message: "password changed",
	}, nil
}
