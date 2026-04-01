package resetpassword

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
)

func TestUseCaseExecute(t *testing.T) {
	userID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	tests := []struct {
		name    string
		setup   func(*MockRepository)
		wantErr error
	}{
		{
			name: "propagates recovery not found when recovery lookup fails",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByRecovery("user@example.com", "123456", gomock.Any()).
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrUserRecoveryNotFound,
		},
		{
			name: "returns update password error when password update fails",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByRecovery("user@example.com", "123456", gomock.Any()).
					Return(entities.User{ID: userID}, nil)
				repo.EXPECT().
					UpdatePassword(userID, gomock.Any()).
					Return(errors.New("update failed"))
			},
			wantErr: autherrs.NewErrUpdatePassword(errors.New("update failed")),
		},
		{
			name: "updates password and clears recovery code on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByRecovery("user@example.com", "123456", gomock.Any()).
					Return(entities.User{ID: userID}, nil)
				repo.EXPECT().
					UpdatePassword(userID, gomock.Any()).
					DoAndReturn(func(_ entities.HexID, newPassword string) error {
						if newPassword == "" || newPassword == "new-password" {
							t.Fatal("UpdatePassword() received unhashed password")
						}
						return nil
					})
				repo.EXPECT().
					SetRecoveryCode(userID, "", time.Time{}).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			gotErr := New(repo).Execute(
				authvalues.UserBasic{Email: "user@example.com", Password: "new-password"},
				"123456",
			)
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
