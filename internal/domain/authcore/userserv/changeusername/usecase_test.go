package changeusername

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestUseCaseExecute(t *testing.T) {
	userID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	currentUser := newChangeUsernameUser(t, userID, "old-name", "current-pass")

	tests := []struct {
		name            string
		setup           func(*MockUserRepository, *MockRepository)
		currentPassword string
		wantErr         error
	}{
		{
			name: "propagates wrong password from session check",
			setup: func(_ *MockUserRepository, repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
			},
			currentPassword: "wrong-pass",
			wantErr:         autherrs.ErrWrongPassword,
		},
		{
			name: "returns username lookup error on infra failure",
			setup: func(userRepo *MockUserRepository, repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByUsername("new-name").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.NewErrChangeUsernameLookup(errors.New("boom")),
		},
		{
			name: "returns username in use when lookup finds another user",
			setup: func(userRepo *MockUserRepository, repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByUsername("new-name").
					Return(entities.User{ID: userID}, nil)
			},
			wantErr: autherrs.ErrUsernameInUse,
		},
		{
			name: "returns update username error when update fails",
			setup: func(userRepo *MockUserRepository, repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByUsername("new-name").
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", "new-name", errors.New("missing")))
				repo.EXPECT().UpdateUsername(userID, "new-name").Return(errors.New("update failed"))
			},
			wantErr: autherrs.NewErrUpdateUsername(errors.New("update failed")),
		},
		{
			name: "updates username on success",
			setup: func(userRepo *MockUserRepository, repo *MockRepository) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByUsername("new-name").
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", "new-name", errors.New("missing")))
				repo.EXPECT().UpdateUsername(userID, "new-name").Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := NewMockUserRepository(ctrl)
			repo := NewMockRepository(ctrl)
			tt.setup(userRepo, repo)

			password := tt.currentPassword
			if password == "" {
				password = "current-pass"
			}

			gotErr := New(userRepo, repo).Execute("auth-token", authvalues.UserBasic{
				Username: "new-name",
				Password: password,
			})
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

func newChangeUsernameUser(
	t *testing.T,
	id entities.HexID,
	username string,
	password string,
) entities.User {
	t.Helper()

	encryptedPassword, err := authvalues.UserBasic{Password: password}.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}

	return entities.User{
		ID: id,
		UserBasic: authvalues.UserBasic{
			Username: username,
			Password: string(encryptedPassword),
		},
	}
}
