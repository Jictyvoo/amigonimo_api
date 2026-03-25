package login

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestUseCaseExecute(t *testing.T) {
	userID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	validPassword := "s3cret-pass"
	userWithPassword := newTestUser(t, userID, "user@example.com", "amigo", validPassword)
	lookupErr := errors.New("lookup failed")
	tokenLookupErr := errors.New("token lookup failed")
	upsertErr := errors.New("upsert failed")

	tests := []struct {
		name        string
		formUser    entities.UserBasic
		setup       func(*MockUserRepository, *MockTokenRepository)
		wantErr     error
		assertToken func(*testing.T, entities.AuthenticationToken)
	}{
		{
			name:     "returns invalid credentials when login fields are empty",
			formUser: entities.UserBasic{Password: validPassword},
			wantErr:  autherrs.ErrInvalidCredentials,
		},
		{
			name: "returns auth login error on user lookup failure",
			formUser: entities.UserBasic{
				Email:    "user@example.com",
				Password: validPassword,
			},
			setup: func(userRepo *MockUserRepository, _ *MockTokenRepository) {
				userRepo.EXPECT().
					GetUserByEmailOrUsername("user@example.com", "").
					Return(entities.User{}, lookupErr)
			},
			wantErr: autherrs.NewErrLogin(lookupErr),
		},
		{
			name: "returns invalid credentials when user is not found",
			formUser: entities.UserBasic{
				Email:    "user@example.com",
				Password: validPassword,
			},
			setup: func(userRepo *MockUserRepository, _ *MockTokenRepository) {
				userRepo.EXPECT().
					GetUserByEmailOrUsername("user@example.com", "").
					Return(
						entities.User{},
						dberrs.NewErrDatabaseNotFound("user", "user@example.com", errors.New("missing")),
					)
			},
			wantErr: autherrs.ErrInvalidCredentials,
		},
		{
			name: "returns invalid credentials when password does not match",
			formUser: entities.UserBasic{
				Email:    "user@example.com",
				Password: "wrong-pass",
			},
			setup: func(userRepo *MockUserRepository, _ *MockTokenRepository) {
				userRepo.EXPECT().
					GetUserByEmailOrUsername("user@example.com", "").
					Return(userWithPassword, nil)
			},
			wantErr: autherrs.ErrInvalidCredentials,
		},
		{
			name: "returns token lookup error when token fetch fails",
			formUser: entities.UserBasic{
				Email:    "user@example.com",
				Password: validPassword,
			},
			setup: func(userRepo *MockUserRepository, tokenRepo *MockTokenRepository) {
				userRepo.EXPECT().
					GetUserByEmailOrUsername("user@example.com", "").
					Return(userWithPassword, nil)
				tokenRepo.EXPECT().
					GetAuthenticationToken(userID).
					Return(entities.AuthenticationToken{}, tokenLookupErr)
			},
			wantErr: autherrs.NewErrTokenLookup(tokenLookupErr),
		},
		{
			name: "returns update auth token error when upsert fails",
			formUser: entities.UserBasic{
				Email:    "user@example.com",
				Password: validPassword,
			},
			setup: func(userRepo *MockUserRepository, tokenRepo *MockTokenRepository) {
				userRepo.EXPECT().
					GetUserByEmailOrUsername("user@example.com", "").
					Return(userWithPassword, nil)
				tokenRepo.EXPECT().
					GetAuthenticationToken(userID).
					Return(entities.AuthenticationToken{}, dberrs.NewErrDatabaseNotFound("auth_token", userID.String(), errors.New("missing")))
				tokenRepo.EXPECT().
					UpsertAuthToken(gomock.Any()).
					DoAndReturn(func(authentication *entities.AuthenticationToken) error {
						if authentication == nil {
							t.Fatal("UpsertAuthToken() received nil token")
						}
						if authentication.AuthToken == "" {
							t.Fatal("UpsertAuthToken() received empty auth token")
						}
						if !authentication.RefreshToken.Valid {
							t.Fatal("UpsertAuthToken() received invalid refresh token")
						}
						return upsertErr
					})
			},
			wantErr: autherrs.NewErrUpdateAuthToken(upsertErr),
		},
		{
			name: "returns regenerated token with user on success",
			formUser: entities.UserBasic{
				Username: "amigo",
				Password: validPassword,
			},
			setup: func(userRepo *MockUserRepository, tokenRepo *MockTokenRepository) {
				existingToken := entities.AuthenticationToken{
					BasicAuthToken: entities.BasicAuthToken{
						AuthToken: "old-token",
						ExpiresAt: time.Now().Add(-time.Hour),
						RefreshToken: uuid.NullUUID{
							UUID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
							Valid: true,
						},
					},
					ID: userID,
				}

				userRepo.EXPECT().
					GetUserByEmailOrUsername("", "amigo").
					Return(userWithPassword, nil)
				tokenRepo.EXPECT().
					GetAuthenticationToken(userID).
					Return(existingToken, nil)
				tokenRepo.EXPECT().
					UpsertAuthToken(gomock.Any()).
					DoAndReturn(func(authentication *entities.AuthenticationToken) error {
						if authentication == nil {
							t.Fatal("UpsertAuthToken() received nil token")
						}
						if authentication.AuthToken == "" ||
							authentication.AuthToken == "old-token" {
							t.Fatal("UpsertAuthToken() did not receive a regenerated auth token")
						}
						if !authentication.RefreshToken.Valid {
							t.Fatal("UpsertAuthToken() received invalid refresh token")
						}
						return nil
					})
			},
			assertToken: func(t *testing.T, authToken entities.AuthenticationToken) {
				t.Helper()

				if authToken.AuthToken == "" || authToken.AuthToken == "old-token" {
					t.Fatalf("AuthToken = %q, want regenerated token", authToken.AuthToken)
				}
				if !authToken.RefreshToken.Valid {
					t.Fatal("RefreshToken.Valid = false, want true")
				}
				if authToken.User.ID != userID {
					t.Fatalf("User.ID = %s, want %s", authToken.User.ID, userID)
				}
				if authToken.User.Email != userWithPassword.Email {
					t.Fatalf(
						"User.Email = %q, want %q",
						authToken.User.Email,
						userWithPassword.Email,
					)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := NewMockUserRepository(ctrl)
			tokenRepo := NewMockTokenRepository(ctrl)

			if tt.setup != nil {
				tt.setup(userRepo, tokenRepo)
			}

			useCase := New(userRepo, tokenRepo)
			gotToken, gotErr := useCase.Execute(tt.formUser)

			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("Execute() error = %v, want nil", gotErr)
			}

			if tt.assertToken != nil {
				tt.assertToken(t, gotToken)
			}
		})
	}
}

func newTestUser(
	t *testing.T,
	id entities.HexID,
	email string,
	username string,
	password string,
) entities.User {
	t.Helper()

	encryptedPassword, err := entities.UserBasic{Password: password}.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}

	return entities.User{
		UserBasic: entities.UserBasic{
			Email:    email,
			Username: username,
			Password: string(encryptedPassword),
		},
		ID: id,
	}
}
