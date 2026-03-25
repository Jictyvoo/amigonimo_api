package lookuprecovery

import (
	"errors"
	"testing"

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
	lookupErr := errors.New("lookup failed")
	user := entities.User{
		ID: userID,
		UserBasic: entities.UserBasic{
			Email: "user@example.com",
		},
	}

	tests := []struct {
		name     string
		setup    func(*MockRepository)
		wantMail string
		wantErr  error
	}{
		{
			name: "returns recovery lookup error on infra failure",
			setup: func(repo *MockRepository) {
				repo.EXPECT().GetUserByUsername("amigo").Return(entities.User{}, lookupErr)
			},
			wantErr: autherrs.NewErrRecoveryLookup(lookupErr),
		},
		{
			name: "returns user not found when repository returns database not found",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetUserByUsername("amigo").
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", "amigo", errors.New("missing")))
			},
			wantErr: autherrs.ErrUserNotFound,
		},
		{
			name: "returns obfuscated email on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().GetUserByUsername("amigo").Return(user, nil)
			},
			wantMail: user.ObfuscateEmail(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			gotMail, gotErr := New(repo).Execute("amigo")
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("Execute() error = %v, want nil", gotErr)
			}
			if gotMail != tt.wantMail {
				t.Fatalf("Execute() mail = %q, want %q", gotMail, tt.wantMail)
			}
		})
	}
}
