package regenerate

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func TestUseCaseExecute(t *testing.T) {
	userID, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}

	upsertErr := errors.New("upsert failed")

	tests := []struct {
		name        string
		setup       func(*MockRepository)
		wantErr     error
		assertToken func(*testing.T, entities.AuthenticationToken)
	}{
		{
			name: "returns invalid auth token when repository errors",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					CheckAuthenticationByRefreshToken("refresh").
					Return(entities.AuthenticationToken{}, errors.New("boom"))
			},
			wantErr: autherrs.ErrInvalidAuthToken,
		},
		{
			name: "returns invalid auth token when authentication is expired",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					CheckAuthenticationByRefreshToken("refresh").
					Return(entities.AuthenticationToken{
						BasicAuthToken: entities.BasicAuthToken{
							ExpiresAt: time.Now().Add(-time.Minute),
						},
						User: entities.User{ID: userID},
					}, nil)
			},
			wantErr: autherrs.ErrInvalidAuthToken,
		},
		{
			name: "returns update auth token error when upsert fails",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					CheckAuthenticationByRefreshToken("refresh").
					Return(entities.AuthenticationToken{
						BasicAuthToken: entities.BasicAuthToken{
							AuthToken: "old-token",
							ExpiresAt: time.Now().Add(time.Minute),
							RefreshToken: uuid.NullUUID{
								UUID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
								Valid: true,
							},
						},
						User: entities.User{ID: userID},
					}, nil)
				repo.EXPECT().
					UpsertAuthToken(gomock.Any()).
					Return(upsertErr)
			},
			wantErr: autherrs.NewErrUpdateAuthToken(upsertErr),
		},
		{
			name: "returns regenerated token on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					CheckAuthenticationByRefreshToken("refresh").
					Return(entities.AuthenticationToken{
						BasicAuthToken: entities.BasicAuthToken{
							AuthToken: "old-token",
							ExpiresAt: time.Now().Add(time.Minute),
							RefreshToken: uuid.NullUUID{
								UUID:  uuid.MustParse("11111111-1111-1111-1111-111111111111"),
								Valid: true,
							},
						},
						User: entities.User{ID: userID},
					}, nil)
				repo.EXPECT().
					UpsertAuthToken(gomock.Any()).
					DoAndReturn(func(token *entities.AuthenticationToken) error {
						if token == nil || token.AuthToken == "" || token.AuthToken == "old-token" {
							t.Fatal("UpsertAuthToken() received non-regenerated token")
						}
						return nil
					})
			},
			assertToken: func(t *testing.T, token entities.AuthenticationToken) {
				t.Helper()
				if token.AuthToken == "" || token.AuthToken == "old-token" {
					t.Fatalf("AuthToken = %q, want regenerated token", token.AuthToken)
				}
				if !token.RefreshToken.Valid {
					t.Fatal("RefreshToken.Valid = false, want true")
				}
				if token.User.ID != userID {
					t.Fatalf("User.ID = %s, want %s", token.User.ID, userID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			gotToken, gotErr := New(repo).Execute("refresh")
			if tt.wantErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Fatalf("Execute() error = %v, want %v", gotErr, tt.wantErr)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("Execute() error = %v, want nil", gotErr)
			}
			tt.assertToken(t, gotToken)
		})
	}
}
