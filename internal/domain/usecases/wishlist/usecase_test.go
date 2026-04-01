package wishlist

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type participantFacadeAdapter struct {
	ports.Facade
	*MockparticipantFacadePort
}

func requireAppErrorCode(t *testing.T, err error, wantCode string) {
	t.Helper()
	appErr, ok := errors.AsType[apperr.Contract](err)
	if !ok {
		t.Fatalf("error type = %T, want apperr.Contract", err)
	}
	if appErr.Code() != wantCode {
		t.Fatalf("Code() = %q, want %q", appErr.Code(), wantCode)
	}
}

func TestWishlistUseCase(t *testing.T) {
	sfID := mustHexID(t)
	userID := mustHexID(t)
	participantID := mustHexID(t)
	itemID := mustHexID(t)
	associatedUser := entities.User{ID: userID}
	ref := ParticipantRef{ParticipantID: participantID, UserID: userID, SecretFriendID: sfID}

	t.Run(
		"get wishlist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			fac := NewMockparticipantFacadePort(ctrl)
			repo.EXPECT().
				GetWishlistByParticipant(ParticipantRef{UserID: userID, SecretFriendID: sfID}).
				Return([]WishlistItem{{ID: itemID}}, nil)
			uc := New(
				associatedUser,
				repo,
				participantFacadeAdapter{MockparticipantFacadePort: fac},
			)
			items, err := uc.GetWishlist(sfID)
			if err != nil || len(items) != 1 {
				t.Fatalf("GetWishlist() = (%v, %v), want one item and nil", items, err)
			}
		},
	)

	t.Run(
		"add item", func(t *testing.T) {
			tests := []struct {
				name    string
				setup   func(*MockRepository, *MockparticipantFacadePort)
				wantErr error
			}{
				{
					name: "forbidden when user is not participant",
					setup: func(_ *MockRepository, fac *MockparticipantFacadePort) {
						fac.EXPECT().
							CheckParticipantInSecretFriend(sfID, userID).
							Return(entities.Participant{}, errors.New("missing"))
					},
					wantErr: errors.New("missing"),
				},
				{
					name: "conflict when wishlist capacity reached",
					setup: func(repo *MockRepository, fac *MockparticipantFacadePort) {
						fac.EXPECT().
							CheckParticipantInSecretFriend(sfID, userID).
							Return(entities.Participant{ID: participantID}, nil)
						repo.EXPECT().
							GetWishlistByParticipant(ref).
							Return(make([]WishlistItem, 10), nil)
					},
					wantErr: errors.New("capacity"),
				},
				{
					name: "success",
					setup: func(repo *MockRepository, fac *MockparticipantFacadePort) {
						fac.EXPECT().
							CheckParticipantInSecretFriend(sfID, userID).
							Return(entities.Participant{ID: participantID}, nil)
						repo.EXPECT().GetWishlistByParticipant(ref).Return(nil, nil)
						repo.EXPECT().
							AddWishlistItem(ref, gomock.Any()).
							Return(WishlistItem{ID: itemID}, nil)
					},
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						ctrl := gomock.NewController(t)
						repo := NewMockRepository(ctrl)
						fac := NewMockparticipantFacadePort(ctrl)
						tt.setup(repo, fac)
						uc := New(
							associatedUser,
							repo,
							participantFacadeAdapter{MockparticipantFacadePort: fac},
						)
						_, err := uc.AddItem(sfID, "book", "great")
						if tt.wantErr != nil {
							wantCode := "wishlist_access_forbidden"
							if tt.name == "conflict when wishlist capacity reached" {
								wantCode = "wishlist_capacity_reached"
							}
							requireAppErrorCode(t, err, wantCode)
							return
						}
						if err != nil {
							t.Fatalf("AddItem() error = %v, want nil", err)
						}
					},
				)
			}
		},
	)

	t.Run(
		"delete item", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			fac := NewMockparticipantFacadePort(ctrl)
			fac.EXPECT().
				CheckParticipantInSecretFriend(sfID, userID).
				Return(entities.Participant{ID: participantID}, nil)
			repo.EXPECT().RemoveWishlistItem(itemID, ref).Return(nil)
			uc := New(
				associatedUser,
				repo,
				participantFacadeAdapter{MockparticipantFacadePort: fac},
			)
			if err := uc.DeleteItem(sfID, itemID); err != nil {
				t.Fatalf("DeleteItem() error = %v, want nil", err)
			}
		},
	)
}

func mustHexID(t *testing.T) entities.HexID {
	t.Helper()
	id, err := entities.NewHexID()
	if err != nil {
		t.Fatalf("entities.NewHexID() error = %v", err)
	}
	return id
}
