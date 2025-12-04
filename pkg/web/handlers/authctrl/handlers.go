package authctrl

import (
	"github.com/go-fuego/fuego"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type AuthHandlers struct {
	authUC *services.AuthService
}

func NewAuthHandlers(authUC *services.AuthService) *AuthHandlers {
	return &AuthHandlers{authUC: authUC}
}

// Login handles POST /auth/login
func (h *AuthHandlers) Login(
	c fuego.ContextWithBody[LoginRequest],
) (*LoginResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// Convert API model to entity DTO
	entityReq := entities.LoginRequest{
		Email: req.Email,
	}

	entityResp, err := h.authUC.Login(entityReq)
	if err != nil {
		return nil, err
	}

	// Convert entity response to API model
	return &LoginResponse{
		UserID: entityResp.UserID.String(),
		Token:  entityResp.Token,
	}, nil
}

// Register handles POST /auth/register
func (h *AuthHandlers) Register(
	c fuego.ContextWithBody[RegisterRequest],
) (*RegisterResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	// Convert API model to entity DTO
	entityReq := entities.RegisterRequest{
		FullName:   req.FullName,
		Email:      req.Email,
		InviteCode: req.InviteCode,
	}

	entityResp, err := h.authUC.Register(entityReq)
	if err != nil {
		return nil, err
	}

	// Convert entity response to API model
	return &RegisterResponse{
		UserID: entityResp.UserID.String(),
		Token:  entityResp.Token,
	}, nil
}
