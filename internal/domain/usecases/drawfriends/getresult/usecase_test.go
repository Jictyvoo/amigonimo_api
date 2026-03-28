package getresult

import (
	"errors"
	"testing"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type repoStub struct {
	getDrawResultForUser func(secretFriendID, userID entities.HexID) (entities.DrawResultItem, error)
}

func (r repoStub) GetDrawResultForUser(
	secretFriendID, userID entities.HexID,
) (entities.DrawResultItem, error) {
	return r.getDrawResultForUser(secretFriendID, userID)
}

func TestUseCaseExecute(t *testing.T) {
	sfID := mustHexID(t)
	userID := mustHexID(t)

	t.Run("maps repository errors", func(t *testing.T) {
		uc := New(entities.User{ID: userID}, repoStub{
			getDrawResultForUser: func(secretFriendID, resolvedUserID entities.HexID) (entities.DrawResultItem, error) {
				if resolvedUserID != userID {
					t.Fatalf("resolved user = %s, want %s", resolvedUserID, userID)
				}
				return entities.DrawResultItem{}, errors.New("missing")
			},
		})

		if _, err := uc.Execute(Input{SecretFriendID: sfID}); err == nil {
			t.Fatal("Execute() error = nil, want error")
		}
	})

	t.Run("returns repository value", func(t *testing.T) {
		expected := entities.DrawResultItem{
			Giver: entities.Participant{RelatedUser: entities.User{ID: userID}},
		}
		uc := New(entities.User{ID: userID}, repoStub{
			getDrawResultForUser: func(secretFriendID, resolvedUserID entities.HexID) (entities.DrawResultItem, error) {
				if resolvedUserID != userID {
					t.Fatalf("resolved user = %s, want %s", resolvedUserID, userID)
				}
				return expected, nil
			},
		})

		got, err := uc.Execute(Input{SecretFriendID: sfID})
		if err != nil {
			t.Fatalf("Execute() error = %v, want nil", err)
		}
		if got.Giver.RelatedUser.ID != expected.Giver.RelatedUser.ID {
			t.Fatalf("got giver user = %s, want %s", got.Giver.RelatedUser.ID, expected.Giver.RelatedUser.ID)
		}
	})
}

func mustHexID(t *testing.T) entities.HexID {
	t.Helper()
	id, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}
	return id
}
