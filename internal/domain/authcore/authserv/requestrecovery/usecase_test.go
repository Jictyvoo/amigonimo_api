package requestrecovery

import (
	"errors"
	"testing"

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
		setup   func(*MockUserRepository, *MockMailer)
		wantErr error
	}{
		{
			name: "returns user email not found when repository errors",
			setup: func(userRepo *MockUserRepository, _ *MockMailer) {
				userRepo.EXPECT().
					GetUserByEmail("user@example.com").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrUserEmailNotFound,
		},
		{
			name: "returns user email not found when user is empty",
			setup: func(userRepo *MockUserRepository, _ *MockMailer) {
				userRepo.EXPECT().GetUserByEmail("user@example.com").Return(entities.User{}, nil)
			},
			wantErr: autherrs.ErrUserEmailNotFound,
		},
		{
			name: "returns generate recovery code error when persisting code fails",
			setup: func(userRepo *MockUserRepository, _ *MockMailer) {
				userRepo.EXPECT().
					GetUserByEmail("user@example.com").
					Return(entities.User{ID: userID, UserBasic: authvalues.UserBasic{Email: "user@example.com"}}, nil)
				userRepo.EXPECT().
					SetRecoveryCode(userID, gomock.Any(), gomock.Any()).
					Return(errors.New("persist failed"))
			},
			wantErr: autherrs.ErrGenRecoveryCode,
		},
		{
			name: "sends recovery email on success",
			setup: func(userRepo *MockUserRepository, mailer *MockMailer) {
				userRepo.EXPECT().
					GetUserByEmail("user@example.com").
					Return(entities.User{ID: userID, UserBasic: authvalues.UserBasic{Email: "user@example.com"}}, nil)
				userRepo.EXPECT().
					SetRecoveryCode(userID, gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ entities.HexID, code string, _ any) error {
						if len(code) != 11 {
							t.Fatalf("recovery code length = %d, want 11", len(code))
						}
						return nil
					})
				mailer.EXPECT().
					SendPasswordRecoveryEmail("user@example.com", gomock.Any()).
					Do(func(_ string, code string) {
						if len(code) != 11 {
							t.Fatalf("mailer recovery code length = %d, want 11", len(code))
						}
					})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := NewMockUserRepository(ctrl)
			mailer := NewMockMailer(ctrl)
			tt.setup(userRepo, mailer)

			gotErr := New(userRepo, mailer).Execute("user@example.com")
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
