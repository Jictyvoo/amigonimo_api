package secretfriend

import (
	"errors"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

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

func TestUseCaseCreate(t *testing.T) {
	ownerID := mustHexID(t)
	associatedUser := entities.User{ID: ownerID}
	createErr := errors.New("create failed")
	input := CreateInput{
		Name:            "Amigo Secreto",
		Datetime:        time.Now().Add(24 * time.Hour).UTC(),
		Location:        "Salvador",
		MaxDenyListSize: 2,
	}

	tests := []struct {
		name      string
		setup     func(*MockRepository)
		wantErr   error
		assertRes func(*testing.T, entities.SecretFriend)
	}{
		{
			name: "returns create failed when repository errors",
			setup: func(repo *MockRepository) {
				repo.EXPECT().CreateSecretFriend(gomock.Any()).Return(createErr)
			},
			wantErr: createErr,
		},
		{
			name: "creates secret friend on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					CreateSecretFriend(gomock.Any()).
					DoAndReturn(func(sf *entities.SecretFriend) error {
						if sf == nil || sf.ID.IsEmpty() {
							t.Fatal("CreateSecretFriend() received empty secret friend")
						}
						if sf.OwnerID != ownerID {
							t.Fatalf("OwnerID = %s, want %s", sf.OwnerID, ownerID)
						}
						return nil
					})
			},
			assertRes: func(t *testing.T, sf entities.SecretFriend) {
				t.Helper()
				if sf.ID.IsEmpty() {
					t.Fatal("ID is empty")
				}
				if sf.Status != entities.StatusDraft {
					t.Fatalf("Status = %s, want %s", sf.Status, entities.StatusDraft)
				}
				if sf.InviteCode == "" {
					t.Fatal("InviteCode is empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)

			uc := New(associatedUser, repo)
			got, err := uc.Create(input)
			if tt.wantErr != nil {
				requireAppErrorCode(t, err, "secret_friend_create_failed")
				return
			}
			if err != nil {
				t.Fatalf("Create() error = %v, want nil", err)
			}
			tt.assertRes(t, got)
		})
	}
}

func TestUseCaseGetInviteInfo(t *testing.T) {
	code := "invite01"
	inviteErr := errors.New("invite not found")
	want := entities.SecretFriend{ID: mustHexID(t), InviteCode: code}

	tests := []struct {
		name    string
		setup   func(*MockRepository)
		wantErr error
		want    entities.SecretFriend
	}{
		{
			name: "returns invite not found on repository error",
			setup: func(repo *MockRepository) {
				repo.EXPECT().
					GetSecretFriendByInviteCode(code).
					Return(entities.SecretFriend{}, inviteErr)
			},
			wantErr: inviteErr,
		},
		{
			name: "returns secret friend on success",
			setup: func(repo *MockRepository) {
				repo.EXPECT().GetSecretFriendByInviteCode(code).Return(want, nil)
			},
			want: want,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)
			uc := New(entities.User{}, repo)

			got, err := uc.GetInviteInfo(code)
			if tt.wantErr != nil {
				requireAppErrorCode(t, err, "secret_friend_invite_not_found")
				return
			}
			if err != nil {
				t.Fatalf("GetInviteInfo() error = %v, want nil", err)
			}
			if got.ID != tt.want.ID {
				t.Fatalf("GetInviteInfo() id = %s, want %s", got.ID, tt.want.ID)
			}
		})
	}
}

func TestUseCaseListUserSecretFriends(t *testing.T) {
	userID := mustHexID(t)
	listErr := errors.New("list failed")
	now := time.Now()
	rawList := []entities.SecretFriend{
		{
			ID:       mustHexID(t),
			OwnerID:  userID,
			Status:   entities.StatusDraft,
			Datetime: now.Add(24 * time.Hour),
		},
		{
			ID:       mustHexID(t),
			OwnerID:  entities.HexID{},
			Status:   entities.StatusOpen,
			Datetime: now.Add(24 * time.Hour),
		},
		{
			ID:       mustHexID(t),
			OwnerID:  userID,
			Status:   entities.StatusClosed,
			Datetime: now.Add(-24 * time.Hour),
		},
		{
			ID:       mustHexID(t),
			OwnerID:  entities.HexID{},
			Status:   entities.StatusClosed,
			Datetime: now.Add(-24 * time.Hour),
		},
	}

	tests := []struct {
		name      string
		setup     func(*MockRepository)
		wantErr   error
		assertRes func(*testing.T, ActiveInactiveListEvents)
	}{
		{
			name: "returns list failed on repository error",
			setup: func(repo *MockRepository) {
				repo.EXPECT().ListSecretFriends(userID).Return(nil, listErr)
			},
			wantErr: listErr,
		},
		{
			name: "sorts secret friends by activity and ownership",
			setup: func(repo *MockRepository) {
				repo.EXPECT().ListSecretFriends(userID).Return(rawList, nil)
			},
			assertRes: func(t *testing.T, got ActiveInactiveListEvents) {
				t.Helper()
				if len(got.Active.Created) != 1 || len(got.Active.Participant) != 1 ||
					len(got.Inactive.Created) != 1 || len(got.Inactive.Participant) != 1 {
					t.Fatalf("unexpected grouping sizes: %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			tt.setup(repo)
			uc := New(entities.User{ID: userID}, repo)

			got, err := uc.ListUserSecretFriends(userID)
			if tt.wantErr != nil {
				requireAppErrorCode(t, err, "secret_friend_list_failed")
				return
			}
			if err != nil {
				t.Fatalf("ListUserSecretFriends() error = %v, want nil", err)
			}
			tt.assertRes(t, got)
		})
	}
}

func TestUseCaseGetAndOwnershipAndUpdate(t *testing.T) {
	ownerID := mustHexID(t)
	otherID := mustHexID(t)
	sfID := mustHexID(t)
	base := entities.SecretFriend{ID: sfID, OwnerID: ownerID, Name: "old", Location: "old place"}
	getErr := errors.New("not found")
	updateErr := errors.New("update failed")

	t.Run("get returns mapped error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().GetSecretFriendByID(sfID).Return(entities.SecretFriend{}, getErr)
		uc := New(entities.User{ID: ownerID}, repo)

		_, err := uc.Get(sfID)
		requireAppErrorCode(t, err, "secret_friend_not_found")
	})

	t.Run("check owner returns true when owner matches", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().GetSecretFriendByID(sfID).Return(base, nil)
		uc := New(entities.User{ID: ownerID}, repo)
		got, err := uc.CheckUserIsOwner(sfID)
		if err != nil || !got {
			t.Fatalf("CheckUserIsOwner() = (%v, %v), want (true, nil)", got, err)
		}
	})

	t.Run("update forbids non-owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().GetSecretFriendByID(sfID).Return(base, nil)
		uc := New(entities.User{ID: otherID}, repo)
		err := uc.Update(UpdateInput{ID: sfID, Name: "new"})
		if err == nil {
			t.Fatal("Update() error = nil, want forbidden")
		}
	})

	t.Run("update maps repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().GetSecretFriendByID(sfID).Return(base, nil)
		repo.EXPECT().UpdateSecretFriend(gomock.Any()).Return(updateErr)
		uc := New(entities.User{ID: ownerID}, repo)
		err := uc.Update(UpdateInput{ID: sfID, Name: "new"})
		requireAppErrorCode(t, err, "secret_friend_update_failed")
	})

	t.Run("update applies fields on success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		repo.EXPECT().GetSecretFriendByID(sfID).Return(base, nil)
		repo.EXPECT().
			UpdateSecretFriend(gomock.Any()).
			DoAndReturn(func(sf *entities.SecretFriend) error {
				if sf.Name != "new" || sf.Location != "new place" {
					t.Fatalf("UpdateSecretFriend() received %+v", sf)
				}
				return nil
			})
		uc := New(entities.User{ID: ownerID}, repo)
		err := uc.Update(UpdateInput{ID: sfID, Name: "new", Location: "new place"})
		if err != nil {
			t.Fatalf("Update() error = %v, want nil", err)
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
