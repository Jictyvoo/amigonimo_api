package entities

import (
	"testing"
	"time"
)

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

func TestWishlistItemNormalize(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 25, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		item         WishlistItem
		wantLabel    string
		wantComments string
		assertTime   func(*testing.T, WishlistItem)
	}{
		{
			name: "trims label and comments and initializes timestamps",
			item: WishlistItem{
				Label:    "  Board Game  ",
				Comments: "  cooperative  ",
			},
			wantLabel:    "Board Game",
			wantComments: "cooperative",
			assertTime: func(t *testing.T, got WishlistItem) {
				t.Helper()
				if got.CreatedAt.IsZero() {
					t.Fatal("CreatedAt is zero")
				}
				if got.UpdatedAt.IsZero() {
					t.Fatal("UpdatedAt is zero")
				}
			},
		},
		{
			name: "preserves non-zero timestamps",
			item: WishlistItem{
				Label:    " book ",
				Comments: " fun ",
				Timestamp: Timestamp{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantLabel:    "book",
			wantComments: "fun",
			assertTime: func(t *testing.T, got WishlistItem) {
				t.Helper()
				if !got.CreatedAt.Equal(now) {
					t.Fatalf("CreatedAt = %v, want %v", got.CreatedAt, now)
				}
				if !got.UpdatedAt.Equal(now) {
					t.Fatalf("UpdatedAt = %v, want %v", got.UpdatedAt, now)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := tt.item
			item.Normalize()

			if item.Label != tt.wantLabel {
				t.Fatalf("Normalize() Label = %q, want %q", item.Label, tt.wantLabel)
			}
			if item.Comments != tt.wantComments {
				t.Fatalf("Normalize() Comments = %q, want %q", item.Comments, tt.wantComments)
			}
			tt.assertTime(t, item)
		})
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
