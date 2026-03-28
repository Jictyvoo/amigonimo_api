package getresult

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func TestUseCaseExecute(t *testing.T) {
	sfID := mustHexID(t)
	userID := mustHexID(t)

	expected := entities.DrawResultItem{
		Giver: entities.Participant{RelatedUser: entities.User{ID: userID}},
	}

	tests := []struct {
		name     string
		repoItem entities.DrawResultItem
		repoErr  error
		wantErr  bool
		check    func(t *testing.T, got entities.DrawResultItem)
	}{
		{
			name:    "maps repository errors",
			repoErr: errors.New("missing"),
			wantErr: true,
		},
		{
			name:     "returns repository value",
			repoItem: expected,
			check: func(t *testing.T, got entities.DrawResultItem) {
				t.Helper()
				if got.Giver.RelatedUser.ID != expected.Giver.RelatedUser.ID {
					t.Fatalf(
						"got giver user = %s, want %s",
						got.Giver.RelatedUser.ID,
						expected.Giver.RelatedUser.ID,
					)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				repo := NewMockRepository(ctrl)
				repo.EXPECT().
					GetDrawResultForUser(sfID, userID).
					Return(tt.repoItem, tt.repoErr)
				uc := New(entities.User{ID: userID}, repo)

				got, err := uc.Execute(Input{SecretFriendID: sfID})
				if (err != nil) != tt.wantErr {
					t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !tt.wantErr && tt.check != nil {
					tt.check(t, got)
				}
			},
		)
	}
}

func mustHexID(t *testing.T) entities.HexID {
	t.Helper()
	id, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}
	return id
}
