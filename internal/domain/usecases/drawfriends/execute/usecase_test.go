package execute

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type secretFriendFacadeAdapter struct {
	ports.Facade
	*MocksecretFriendFacadePort
}

func TestUseCaseExecute(t *testing.T) {
	mustHexID := func(t *testing.T) entities.HexID {
		t.Helper()
		id, err := entities.NewHexID()
		if err != nil {
			t.Fatalf("entities.NewHexID() error = %v", err)
		}
		return id
	}

	newParticipant := func(t *testing.T) entities.Participant {
		t.Helper()
		return entities.Participant{
			ID:             mustHexID(t),
			SecretFriendID: mustHexID(t),
			RelatedUser: entities.User{
				ID: mustHexID(t),
			},
		}
	}

	sfID := mustHexID(t)
	p1 := newParticipant(t)
	p2 := newParticipant(t)
	p3 := newParticipant(t)
	drawService := New()

	tests := []struct {
		name       string
		setupMocks func(t *testing.T, repo *MockRepository, facade *MocksecretFriendFacadePort)
		wantErr    bool
		check      func(t *testing.T, out Output)
	}{
		{
			name: "returns not found when secret friend lookup fails",
			setupMocks: func(_ *testing.T, _ *MockRepository, facade *MocksecretFriendFacadePort) {
				facade.EXPECT().
					GetSecretFriendByID(sfID).
					Return(entities.SecretFriend{}, errors.New("missing"))
			},
			wantErr: true,
		},
		{
			name: "returns conflict when already drawn",
			setupMocks: func(_ *testing.T, _ *MockRepository, facade *MocksecretFriendFacadePort) {
				facade.EXPECT().
					GetSecretFriendByID(sfID).
					Return(entities.SecretFriend{ID: sfID, Status: entities.StatusDrawn}, nil)
			},
			wantErr: true,
		},
		{
			name: "saves results and updates status on success",
			setupMocks: func(t *testing.T, repo *MockRepository, facade *MocksecretFriendFacadePort) {
				t.Helper()
				facade.EXPECT().
					GetSecretFriendByID(sfID).
					Return(
						entities.SecretFriend{
							ID:           sfID,
							Status:       entities.StatusOpen,
							Participants: []entities.Participant{p1, p2, p3},
						}, nil,
					)
				repo.EXPECT().BeginTx(gomock.Any(), nil).Return(
					func(commit bool) error {
						if !commit {
							t.Fatal("transaction finished without commit")
						}
						return nil
					}, nil,
				)
				repo.EXPECT().
					SaveDrawResults(sfID, gomock.Any()).
					DoAndReturn(
						func(_ entities.HexID, pairs []entities.DrawResultItem) error {
							if len(pairs) != 3 {
								t.Fatalf("SaveDrawResults() pairs = %d, want 3", len(pairs))
							}
							return nil
						},
					)
				facade.EXPECT().UpdateStatus(sfID, entities.StatusDrawn).Return(nil)
			},
			check: func(t *testing.T, out Output) {
				t.Helper()
				if out.ParticipantCount != 3 {
					t.Fatalf("ParticipantCount = %d, want 3", out.ParticipantCount)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				repo := NewMockRepository(ctrl)
				facade := NewMocksecretFriendFacadePort(ctrl)
				tt.setupMocks(t, repo, facade)

				uc := New(
					repo,
					secretFriendFacadeAdapter{MocksecretFriendFacadePort: facade},
					drawService,
				)

				out, err := uc.Execute(Input{SecretFriendID: sfID})
				if (err != nil) != tt.wantErr {
					t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !tt.wantErr && tt.check != nil {
					tt.check(t, out)
				}
			},
		)
	}
}
