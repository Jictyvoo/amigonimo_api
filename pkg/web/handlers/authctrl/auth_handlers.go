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
	c fuego.ContextWithBody[entities.LoginRequest],
) (*entities.LoginResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	resp, err := h.authUC.Login(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Register handles POST /auth/register
func (h *AuthHandlers) Register(
	c fuego.ContextWithBody[entities.RegisterRequest],
) (*entities.RegisterResponse, error) {
	req, err := c.Body()
	if err != nil {
		return nil, err
	}

	resp, err := h.authUC.Register(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
