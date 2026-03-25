package verify

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
		wantErr error
	}{
		{
			name: "returns verification code error when lookup fails",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByVerificationCode("verify-code").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrVerificationCode,
		},
		{
			name: "returns set verification error when update fails",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByVerificationCode("verify-code").
					Return(entities.User{ID: userID}, nil)
				repo.EXPECT().SetUserVerified(userID).Return(errors.New("update failed"))
			},
			wantErr: autherrs.NewErrSetVerification(errors.New("update failed")),
		},
		{
			name: "verifies user on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByVerificationCode("verify-code").
					Return(entities.User{ID: userID}, nil)
				repo.EXPECT().SetUserVerified(userID).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			gotErr := New(repo).Execute("verify-code")
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
