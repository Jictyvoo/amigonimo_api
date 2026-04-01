package session

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

func TestFindAndCheckUser(t *testing.T) {
	userID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	validUser := newSessionUser(t, userID, "s3cret-pass")

	tests := []struct {
		name     string
		setup    func(*MockRepository)
		wantUser entities.User
		wantErr  error
	}{
		{
			name: "returns user not found when repository errors",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByAuthCode("auth-token").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrUserNotFound,
		},
		{
			name: "returns wrong password when password check fails",
			setup: func(repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(validUser, nil)
			},
			wantErr: autherrs.ErrWrongPassword,
		},
		{
			name: "returns user on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(validUser, nil)
			},
			wantUser: validUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			password := "wrong-pass"
			if tt.wantErr == nil {
				password = "s3cret-pass"
			}

			gotUser, gotErr := FindAndCheckUser(repo, "auth-token", password)
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("FindAndCheckUser() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("FindAndCheckUser() error = %v, want nil", gotErr)
			}
			if gotUser.ID != tt.wantUser.ID {
				t.Fatalf("FindAndCheckUser() user ID = %s, want %s", gotUser.ID, tt.wantUser.ID)
			}
		})
	}
}

func newSessionUser(t *testing.T, id entities.HexID, password string) entities.User {
	t.Helper()

	encryptedPassword, err := authvalues.UserBasic{Password: password}.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}

	return entities.User{
		ID: id,
		UserBasic: authvalues.UserBasic{
			Password: string(encryptedPassword),
		},
	}
}
