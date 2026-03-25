package changeemail

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

	currentUser := newChangeEmailUser(t, userID, "amigo", "old@example.com", "current-pass")

	tests := []struct {
		name    string
		setup   func(*MockUserRepository, *MockRepository, *MockMailer)
		form    entities.UserBasic
		wantErr error
	}{
		{
			name: "propagates user not found from session check",
			form: entities.UserBasic{
				Email:    "new@example.com",
				Password: "current-pass",
			},
			setup: func(_ *MockUserRepository, repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().
					GetUserByAuthCode("auth-token").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrUserNotFound,
		},
		{
			name: "returns email in use when new email matches current email",
			form: entities.UserBasic{
				Email:    "old@example.com",
				Password: "current-pass",
			},
			setup: func(_ *MockUserRepository, repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
			},
			wantErr: autherrs.ErrEmailInUse,
		},
		{
			name: "returns change email lookup error on infra failure",
			form: entities.UserBasic{
				Email:    "new@example.com",
				Password: "current-pass",
			},
			setup: func(userRepo *MockUserRepository, repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByEmail("new@example.com").
					Return(entities.User{}, errors.New("boom"))
			},
			wantErr: autherrs.NewErrChangeEmailLookup(errors.New("boom")),
		},
		{
			name: "returns email in use when lookup finds another user",
			form: entities.UserBasic{
				Email:    "new@example.com",
				Password: "current-pass",
			},
			setup: func(userRepo *MockUserRepository, repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByEmail("new@example.com").
					Return(entities.User{ID: userID}, nil)
			},
			wantErr: autherrs.ErrEmailInUse,
		},
		{
			name: "returns update email error when change email fails",
			form: entities.UserBasic{
				Email:    "new@example.com",
				Password: "current-pass",
			},
			setup: func(userRepo *MockUserRepository, repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByEmail("new@example.com").
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", "new@example.com", errors.New("missing")))
				repo.EXPECT().
					ChangeEmail(userID, "new@example.com").
					Return(errors.New("update failed"))
			},
			wantErr: autherrs.NewErrUpdateEmail(errors.New("update failed")),
		},
		{
			name: "returns set verification error when verification update fails",
			form: entities.UserBasic{
				Email:    "new@example.com",
				Password: "current-pass",
			},
			setup: func(userRepo *MockUserRepository, repo *MockRepository, _ *MockMailer) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByEmail("new@example.com").
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", "new@example.com", errors.New("missing")))
				repo.EXPECT().ChangeEmail(userID, "new@example.com").Return(nil)
				repo.EXPECT().
					SetNewVerificationCode(userID, gomock.Any()).
					Return(errors.New("verify failed"))
			},
			wantErr: autherrs.NewErrSetVerification(errors.New("verify failed")),
		},
		{
			name: "changes email and sends activation email on success",
			form: entities.UserBasic{
				Email:    "new@example.com",
				Password: "current-pass",
			},
			setup: func(userRepo *MockUserRepository, repo *MockRepository, mailer *MockMailer) {
				repo.EXPECT().GetUserByAuthCode("auth-token").Return(currentUser, nil)
				userRepo.EXPECT().
					GetUserByEmail("new@example.com").
					Return(entities.User{}, dberrs.NewErrDatabaseNotFound("user", "new@example.com", errors.New("missing")))
				repo.EXPECT().ChangeEmail(userID, "new@example.com").Return(nil)
				repo.EXPECT().
					SetNewVerificationCode(userID, gomock.Any()).
					DoAndReturn(func(_ entities.HexID, token string) error {
						if token == "" {
							t.Fatal("SetNewVerificationCode() received empty token")
						}
						return nil
					})
				mailer.EXPECT().
					SendActivationEmail("new@example.com", gomock.Any()).
					Do(func(_ string, token string) {
						if token == "" {
							t.Fatal("SendActivationEmail() received empty token")
						}
					})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userRepo := NewMockUserRepository(ctrl)
			repo := NewMockRepository(ctrl)
			mailer := NewMockMailer(ctrl)
			tt.setup(userRepo, repo, mailer)

			gotErr := New(userRepo, repo, mailer).Execute("auth-token", tt.form)
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

func newChangeEmailUser(
	t *testing.T,
	id entities.HexID,
	username string,
	email string,
	password string,
) entities.User {
	t.Helper()

	encryptedPassword, err := entities.UserBasic{Password: password}.EncryptPassword()
	if err != nil {
		t.Fatalf("EncryptPassword() error = %v", err)
	}

	return entities.User{
		ID: id,
		UserBasic: entities.UserBasic{
			Username: username,
			Email:    email,
			Password: string(encryptedPassword),
		},
	}
}
