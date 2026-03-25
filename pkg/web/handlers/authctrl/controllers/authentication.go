package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (h *AuthenticationController) SignUp(
	c fuego.Context[FormUser, struct{}],
) (*SuccessResponse, error) {
	req, err := c.Body()
	if err != nil {
		return &SuccessResponse{Message: "Failed to obtain request body"}, h.HandleError(err)
	}

	userDTO := entities.UserBasic{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if _, err = h.authServ.SignUp.Execute(userDTO); err != nil {
		return nil, h.HandleError(err)
	}

	c.SetStatus(http.StatusCreated)
	return &SuccessResponse{Success: true, Message: "User created successfully"}, nil
}

func (h *AuthenticationController) LogIn(
	c fuego.Context[FormUser, struct{}],
) (*LoginResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, h.HandleError(err)
	}

	userDTO := entities.UserBasic{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	authToken, err := h.authServ.LogIn.Execute(userDTO)
	if err != nil {
		return nil, h.HandleError(err)
	}

	// Generate JWT from AuthenticationToken
	jwtToken, err := h.generateJWT(authToken)
	if err != nil {
		return nil, h.HandleError(err)
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

	authToken, err := h.authServ.RegenerateToken.Execute(refreshToken)
	if err != nil {
		return nil, h.HandleError(err)
	}

	// Generate JWT from AuthenticationToken
	jwtToken, err := h.generateJWT(authToken)
	if err != nil {
		return nil, h.HandleError(err)
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
		return nil, h.HTTPError(
			http.StatusPreconditionRequired,
			errors.New("verification code is required"),
		)
	}

	decoded, err := base64.StdEncoding.DecodeString(verifyCode)
	if err != nil {
		return nil, h.HTTPError(
			http.StatusPreconditionRequired,
			errors.New("invalid verification code format"),
		)
	}

	if err = h.authServ.VerifyUser.Execute(string(decoded)); err != nil {
		return nil, h.HandleError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "User successfully verified",
	}, nil
}

func (h *AuthenticationController) ForgotPassword(
	c fuego.Context[FormUser, struct{}],
) (*ForgotPasswordResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, h.HandleError(err)
	}

	if len(req.Username) > 0 {
		obfuscated, lookupErr := h.authServ.LookupRecoveryContact.Execute(req.Username)
		if lookupErr != nil {
			return nil, h.HandleError(lookupErr)
		}
		c.Response().WriteHeader(http.StatusNonAuthoritativeInfo)
		return &ForgotPasswordResponse{
			SuccessResponse: SuccessResponse{Success: true},
			ObfuscatedEmail: obfuscated,
		}, nil
	}

	if len(req.Email) > 0 {
		if err = h.authServ.RequestPasswordRecovery.Execute(req.Email); err != nil {
			return nil, h.HandleError(err)
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

	return nil, h.HTTPError(http.StatusUnprocessableEntity, errors.New("no user or email provided"))
}

func (h *AuthenticationController) CheckRecoveryCode(
	c fuego.Context[FormRecoveryCode, struct{}],
) (*SuccessResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, h.HandleError(err)
	}

	if _, err = h.authServ.CheckRecoveryCode.Execute(req.Email, req.RecoveryCode); err != nil {
		return nil, h.HandleError(err)
	}

	return &SuccessResponse{
		Success: true,
		Message: "sent recovery-code is valid",
	}, nil
}

func (h *AuthenticationController) PasswordReset(
	c fuego.Context[FormResetPassword, struct{}],
) (*SuccessResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, h.HandleError(err)
	}

	userDTO := entities.UserBasic{
		Email:    req.Email,
		Password: req.NewPassword,
	}

	if err = h.authServ.ResetPassword.Execute(userDTO, req.RecoveryCode); err != nil {
		return nil, h.HandleError(err)
	}

	c.Response().WriteHeader(http.StatusCreated)
	return &SuccessResponse{
		Success: true,
		Message: "password changed",
	}, nil
}
