package signup

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock/dberrs"
)

func TestUseCaseExecute(t *testing.T) {
	existingID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	input := entities.UserBasic{
		Username: "amigo",
		Email:    "user@example.com",
		Password: "s3cret-pass",
	}

	tests := []struct {
		name       string
		setup      func(*MockRepository, *MockMailer)
		wantErr    error
		assertUser func(*testing.T, entities.User)
	}{
		{
			name: "returns signup lookup error on infra failure",
			setup: func(repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().
					GetUserByEmailOrUsername(input.Email, input.Username).
					Return(entities.User{}, errors.New("lookup failed"))
			},
			wantErr: autherrs.NewErrSignUpLookup(errors.New("lookup failed")),
		},
		{
			name: "returns email or username used when user already exists",
			setup: func(repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().
					GetUserByEmailOrUsername(input.Email, input.Username).
					Return(entities.User{ID: existingID}, nil)
			},
			wantErr: autherrs.ErrEmailOrUsernameUsed,
		},
		{
			name: "returns user creation error when create user fails",
			setup: func(repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().
					GetUserByEmailOrUsername(input.Email, input.Username).
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", input.Email, errors.New("missing")))
				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(errors.New("create failed"))
			},
			wantErr: autherrs.NewErrUserCreation(errors.New("create failed")),
		},
		{
			name: "creates user and sends activation email on success",
			setup: func(repo *MockRepository, mailer *MockMailer) {
				repo.EXPECT().
					GetUserByEmailOrUsername(input.Email, input.Username).
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", input.Email, errors.New("missing")))
				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					DoAndReturn(func(user entities.User, token string) error {
						if user.ID.IsEmpty() {
							t.Fatal("CreateUser() received empty user ID")
						}
						if user.Password == "" || user.Password == input.Password {
							t.Fatal("CreateUser() received unhashed password")
						}
						if token == "" {
							t.Fatal("CreateUser() received empty verification token")
						}
						return nil
					})
				mailer.EXPECT().
					SendActivationEmail(input.Email, gomock.Any()).
					Do(func(_ string, token string) {
						if token == "" {
							t.Fatal("SendActivationEmail() received empty token")
						}
					})
			},
			assertUser: func(t *testing.T, user entities.User) {
				t.Helper()
				if user.ID.IsEmpty() {
					t.Fatal("user ID is empty")
				}
				if user.Email != input.Email {
					t.Fatalf("user email = %q, want %q", user.Email, input.Email)
				}
				if user.Password == "" || user.Password == input.Password {
					t.Fatal("user password was not encrypted")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			mailer := NewMockMailer(ctrl)
			tt.setup(repo, mailer)

			gotUser, gotErr := New(repo, mailer).Execute(input)
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("Execute() error = %v, want nil", gotErr)
			}
			tt.assertUser(t, gotUser)
		})
	}
}
