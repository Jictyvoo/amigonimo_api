package entities

import "testing"

func TestNewParticipant(t *testing.T) {
	t.Parallel()

	sfID := mustHexID(t)
	userID := mustHexID(t)
	user := User{ID: userID}

	got := NewParticipant(sfID, user)

	if got.SecretFriendID != sfID {
		t.Fatalf("NewParticipant() SecretFriendID = %v, want %v", got.SecretFriendID, sfID)
	}
	if got.RelatedUser != user {
		t.Fatalf("NewParticipant() RelatedUser = %#v, want %#v", got.RelatedUser, user)
	}
}

func mustHexID(t *testing.T) HexID {
	t.Helper()

	id, err := NewHexID()
	if err != nil {
		t.Fatalf("NewHexID() error = %v, want nil", err)
	}

	return id
}
