package drawfriends

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

func TestUseCaseExecuteAndGetResult(t *testing.T) {
	sfID := mustHexID(t)
	p1 := newParticipant(t)
	p2 := newParticipant(t)
	p3 := newParticipant(t)
	drawService := drawserv.New()

	t.Run("execute returns not found when secret friend lookup fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().
			GetSecretFriendByID(sfID).
			Return(entities.SecretFriend{}, errors.New("missing"))
		uc := New(repo, drawService)
		_, err := uc.Execute(ExecuteInput{SecretFriendID: sfID})
		if err == nil {
			t.Fatal("Execute() error = nil, want error")
		}
	})

	t.Run("execute returns conflict when already drawn", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().
			GetSecretFriendByID(sfID).
			Return(entities.SecretFriend{ID: sfID, Status: entities.StatusDrawn}, nil)
		uc := New(repo, drawService)
		_, err := uc.Execute(ExecuteInput{SecretFriendID: sfID})
		if err == nil {
			t.Fatal("Execute() error = nil, want conflict")
		}
	})

	t.Run("execute saves results and updates status on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		sf := entities.SecretFriend{
			ID:           sfID,
			Status:       entities.StatusOpen,
			Participants: []entities.Participant{p1, p2, p3},
		}
		repo.EXPECT().GetSecretFriendByID(sfID).Return(sf, nil)
		repo.EXPECT().BeginTx(gomock.Any(), nil).Return(func(commit bool) error {
			if !commit {
				t.Fatal("transaction finished without commit")
			}
			return nil
		}, nil)
		repo.EXPECT().
			SaveDrawResults(sfID, gomock.Any()).
			DoAndReturn(func(_ entities.HexID, pairs []entities.DrawResultItem) error {
				if len(pairs) != 3 {
					t.Fatalf("SaveDrawResults() pairs = %d, want 3", len(pairs))
				}
				return nil
			})
		repo.EXPECT().
			UpdateSecretFriend(gomock.Any()).
			DoAndReturn(func(updated *entities.SecretFriend) error {
				if updated.Status != entities.StatusDrawn {
					t.Fatalf("updated status = %s, want %s", updated.Status, entities.StatusDrawn)
				}
				return nil
			})
		uc := New(repo, drawService)
		out, err := uc.Execute(ExecuteInput{SecretFriendID: sfID})
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		if out.ParticipantCount != 3 {
			t.Fatalf("ParticipantCount = %d, want 3", out.ParticipantCount)
		}
	})

	t.Run("get result maps repository error and returns value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		uc := New(repo, drawService)

		repo.EXPECT().
			GetDrawResultForUser(sfID, p1.RelatedUser.ID).
			Return(entities.DrawResultItem{}, errors.New("missing"))
		if _, err := uc.GetResult(GetResultInput{SecretFriendID: sfID, UserID: p1.RelatedUser.ID}); err == nil {
			t.Fatal("GetResult() error = nil, want error")
		}
	})
}

func newParticipant(t *testing.T) entities.Participant {
	t.Helper()
	return entities.Participant{
		ID:             mustHexID(t),
		SecretFriendID: mustHexID(t),
		RelatedUser: entities.User{
			ID: mustHexID(t),
		},
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

var (
	_ dbrock.Transactioner = (*MockRepository)(nil)
	_                      = context.Background
)
