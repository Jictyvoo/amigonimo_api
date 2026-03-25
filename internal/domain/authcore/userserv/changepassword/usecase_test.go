package changepassword

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func TestUseCaseExecute(t *testing.T) {
	userID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	validUser := newChangePasswordUser(t, userID, "current-pass")

	tests := []struct {
		name    string
		setup   func(*MockRepository, *MockSessionRepository)
		wantErr error
	}{
		{
			name: "propagates user not found from session check",
			setup: func(_ *MockRepository, sessionRepo *MockSessionRepository) {
				sessionRepo.EXPECT().
					GetUserByAuthCode("auth-token").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrUserNotFound,
		},
		{
			name: "returns update password error when repository update fails",
			setup: func(passwordRepo *MockRepository, sessionRepo *MockSessionRepository) {
				sessionRepo.EXPECT().
					GetUserByAuthCode("auth-token").
					Return(validUser, nil)
				passwordRepo.EXPECT().
					UpdatePassword(userID, gomock.Any()).
					Return(errors.New("update failed"))
			},
			wantErr: autherrs.NewErrUpdatePassword(errors.New("update failed")),
		},
		{
			name: "updates password on success",
			setup: func(passwordRepo *MockRepository, sessionRepo *MockSessionRepository) {
				sessionRepo.EXPECT().
					GetUserByAuthCode("auth-token").
					Return(validUser, nil)
				passwordRepo.EXPECT().
					UpdatePassword(userID, gomock.Any()).
					DoAndReturn(func(_ entities.HexID, password string) error {
						if password == "" || password == "new-pass" {
							t.Fatal("UpdatePassword() received unhashed password")
						}
						return nil
					})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			passwordRepo := NewMockRepository(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			tt.setup(passwordRepo, sessionRepo)

			gotErr := New(
				passwordRepo,
				sessionRepo,
			).Execute("auth-token", "current-pass", "new-pass")
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("Execute() error = %v, want nil", gotErr)
			}
		})
	}
}

func newChangePasswordUser(t *testing.T, id entities.HexID, password string) entities.User {
	t.Helper()

	encryptedPassword, err := entities.UserBasic{Password: password}.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}

	return entities.User{
		ID: id,
		UserBasic: entities.UserBasic{
			Password: string(encryptedPassword),
		},
	}
}
