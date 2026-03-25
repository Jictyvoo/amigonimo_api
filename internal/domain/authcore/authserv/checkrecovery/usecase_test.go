package checkrecovery

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

	tests := []struct {
		name    string
		setup   func(*MockRepository)
		wantID  entities.HexID
		wantErr error
	}{
		{
			name: "returns user recovery not found when repository errors",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByRecovery("user@example.com", "123456", gomock.Any()).
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrUserRecoveryNotFound,
		},
		{
			name: "returns user recovery not found when user is empty",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByRecovery("user@example.com", "123456", gomock.Any()).
					Return(entities.User{}, nil)
			},
			wantErr: autherrs.ErrUserRecoveryNotFound,
		},
		{
			name: "returns user id on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByRecovery("user@example.com", "123456", gomock.Any()).
					Return(entities.User{ID: userID}, nil)
			},
			wantID: userID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			gotID, gotErr := New(repo).Execute("user@example.com", "123456")
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("Execute() error = %v, want nil", gotErr)
			}
			if gotID != tt.wantID {
				t.Fatalf("Execute() id = %s, want %s", gotID, tt.wantID)
			}
		})
	}
}
