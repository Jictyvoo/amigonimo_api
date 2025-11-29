package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate mockgen -destination=../mocks/user_repository_mock.go -package=mocks github.com/jictyvoo/amigonimo_api/internal/domain/services UserRepository

type (
	UserRepository interface {
		GetUserByEmail(email string) (entities.User, error)
		CreateUser(user entities.User) (entities.User, error)
	}
	AuthService struct {
		userRepo UserRepository
	}
)

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login handles user login by email
func (uc *AuthService) Login(req entities.LoginRequest) (*entities.LoginResponse, error) {
	user, err := uc.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.ID.IsEmpty() {
		return nil, errors.New("user not found")
	}

	token, err := uc.generateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entities.LoginResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}

// Register handles user registration
func (uc *AuthService) Register(req entities.RegisterRequest) (*entities.RegisterResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if !existingUser.ID.IsEmpty() {
		// User exists, return login response
		token, err := uc.generateToken(existingUser.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}
		return &entities.RegisterResponse{
			UserID: existingUser.ID,
			Token:  token,
		}, nil
	}

	// Create new user
	user := entities.User{
		ID:       entities.NewHexID(),
		FullName: req.Fullname,
		Email:    req.Email,
	}

	createdUser, err := uc.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := uc.generateToken(createdUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entities.RegisterResponse{
		UserID: createdUser.ID,
		Token:  token,
	}, nil
}

// generateToken generates a simple token (in production, use JWT)
func (uc *AuthService) generateToken(userID entities.HexID) (string, error) {
	// Simple token generation - in production, use proper JWT
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return fmt.Sprintf("%s:%s", userID.String(), token), nil
}
