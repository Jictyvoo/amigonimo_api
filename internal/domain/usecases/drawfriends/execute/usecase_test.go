package execute

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

type repoStub struct {
	getSecretFriendByID func(id entities.HexID) (entities.SecretFriend, error)
	beginTx             func(ctx context.Context, txOpts *sql.TxOptions) (dbrock.OnFinishFunc, error)
	saveDrawResults     func(secretFriendID entities.HexID, results []entities.DrawResultItem) error
	updateSecretFriend  func(sf *entities.SecretFriend) error
}

func (r repoStub) BeginTx(ctx context.Context, txOpts *sql.TxOptions) (dbrock.OnFinishFunc, error) {
	return r.beginTx(ctx, txOpts)
}

func (r repoStub) GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error) {
	return r.getSecretFriendByID(id)
}

func (r repoStub) UpdateSecretFriend(sf *entities.SecretFriend) error {
	return r.updateSecretFriend(sf)
}

func (r repoStub) SaveDrawResults(secretFriendID entities.HexID, results []entities.DrawResultItem) error {
	return r.saveDrawResults(secretFriendID, results)
}

func TestUseCaseExecute(t *testing.T) {
	sfID := mustHexID(t)
	p1 := newParticipant(t)
	p2 := newParticipant(t)
	p3 := newParticipant(t)
	drawService := drawserv.New()

	t.Run("returns not found when secret friend lookup fails", func(t *testing.T) {
		uc := New(repoStub{
			getSecretFriendByID: func(id entities.HexID) (entities.SecretFriend, error) {
				return entities.SecretFriend{}, errors.New("missing")
			},
		}, drawService)

		if _, err := uc.Execute(Input{SecretFriendID: sfID}); err == nil {
			t.Fatal("Execute() error = nil, want error")
		}
	})

	t.Run("returns conflict when already drawn", func(t *testing.T) {
		uc := New(repoStub{
			getSecretFriendByID: func(id entities.HexID) (entities.SecretFriend, error) {
				return entities.SecretFriend{ID: sfID, Status: entities.StatusDrawn}, nil
			},
		}, drawService)

		if _, err := uc.Execute(Input{SecretFriendID: sfID}); err == nil {
			t.Fatal("Execute() error = nil, want conflict")
		}
	})

	t.Run("saves results and updates status on success", func(t *testing.T) {
		uc := New(repoStub{
			getSecretFriendByID: func(id entities.HexID) (entities.SecretFriend, error) {
				return entities.SecretFriend{
					ID:           sfID,
					Status:       entities.StatusOpen,
					Participants: []entities.Participant{p1, p2, p3},
				}, nil
			},
			beginTx: func(ctx context.Context, txOpts *sql.TxOptions) (dbrock.OnFinishFunc, error) {
				return func(commit bool) error {
					if !commit {
						t.Fatal("transaction finished without commit")
					}
					return nil
				}, nil
			},
			saveDrawResults: func(_ entities.HexID, pairs []entities.DrawResultItem) error {
				if len(pairs) != 3 {
					t.Fatalf("SaveDrawResults() pairs = %d, want 3", len(pairs))
				}
				return nil
			},
			updateSecretFriend: func(updated *entities.SecretFriend) error {
				if updated.Status != entities.StatusDrawn {
					t.Fatalf("updated status = %s, want %s", updated.Status, entities.StatusDrawn)
				}
				return nil
			},
		}, drawService)

		out, err := uc.Execute(Input{SecretFriendID: sfID})
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		if out.ParticipantCount != 3 {
			t.Fatalf("ParticipantCount = %d, want 3", out.ParticipantCount)
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
