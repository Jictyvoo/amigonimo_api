package participant

import (
	"errors"
	"testing"

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

func TestParticipantUseCase(t *testing.T) {
	sfID := mustHexID(t)
	userID := mustHexID(t)
	participantID := mustHexID(t)
	associatedUser := entities.User{ID: userID}
	participantEntity := entities.Participant{ID: participantID, RelatedUser: associatedUser}

	t.Run("confirm participation", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func(*MockRepository, *MockSecretFriendFacade)
			wantErr error
		}{
			{
				name: "secret friend not found",
				setup: func(_ *MockRepository, sf *MockSecretFriendFacade) {
					sf.EXPECT().
						GetSecretFriendByID(sfID).
						Return(entities.SecretFriend{}, errors.New("missing"))
				},
				wantErr: errors.New("missing"),
			},
			{
				name: "confirm participation failed",
				setup: func(repo *MockRepository, sf *MockSecretFriendFacade) {
					sf.EXPECT().
						GetSecretFriendByID(sfID).
						Return(entities.SecretFriend{ID: sfID}, nil)
					repo.EXPECT().
						AddParticipant(sfID, userID).
						Return(entities.Participant{}, errors.New("add failed"))
				},
				wantErr: errors.New("add failed"),
			},
			{
				name: "success",
				setup: func(repo *MockRepository, sf *MockSecretFriendFacade) {
					sf.EXPECT().
						GetSecretFriendByID(sfID).
						Return(entities.SecretFriend{ID: sfID}, nil)
					repo.EXPECT().AddParticipant(sfID, userID).Return(participantEntity, nil)
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				repo := NewMockRepository(ctrl)
				sf := NewMockSecretFriendFacade(ctrl)
				tt.setup(repo, sf)
				uc := New(associatedUser, repo, sf)
				_, err := uc.ConfirmParticipation(sfID)
				if tt.wantErr != nil {
					wantCode := "secret_friend_not_found"
					if tt.name == "confirm participation failed" {
						wantCode = "participant_confirm_failed"
					}
					requireAppErrorCode(t, err, wantCode)
					return
				}
				if err != nil {
					t.Fatalf("ConfirmParticipation() error = %v, want nil", err)
				}
			})
		}
	})

	t.Run("list participants", func(t *testing.T) {
		tests := []struct {
			name    string
			setup   func(*MockRepository, *MockSecretFriendFacade)
			wantErr error
		}{
			{
				name: "forbidden for non participant non owner",
				setup: func(repo *MockRepository, sf *MockSecretFriendFacade) {
					repo.EXPECT().
						GetParticipant(sfID, userID).
						Return(entities.Participant{}, errors.New("missing"))
					sf.EXPECT().CheckUserIsOwner(sfID).Return(false, nil)
				},
				wantErr: errors.New("forbidden"),
			},
			{
				name: "owner can list",
				setup: func(repo *MockRepository, sf *MockSecretFriendFacade) {
					repo.EXPECT().
						GetParticipant(sfID, userID).
						Return(entities.Participant{}, errors.New("missing"))
					sf.EXPECT().CheckUserIsOwner(sfID).Return(true, nil)
					repo.EXPECT().
						ListParticipants(sfID).
						Return([]entities.Participant{participantEntity}, nil)
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				repo := NewMockRepository(ctrl)
				sf := NewMockSecretFriendFacade(ctrl)
				tt.setup(repo, sf)
				uc := New(associatedUser, repo, sf)
				_, err := uc.ListParticipants(sfID)
				if tt.wantErr != nil {
					requireAppErrorCode(t, err, "participant_list_forbidden")
					return
				}
				if err != nil {
					t.Fatalf("ListParticipants() error = %v, want nil", err)
				}
			})
		}
	})

	t.Run("mark as ready", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		sf := NewMockSecretFriendFacade(ctrl)
		repo.EXPECT().GetParticipant(sfID, userID).Return(participantEntity, nil)
		repo.EXPECT().SetParticipantReady(sfID, userID, true).Return(nil)
		uc := New(associatedUser, repo, sf)
		if err := uc.MarkAsReady(sfID); err != nil {
			t.Fatalf("MarkAsReady() error = %v, want nil", err)
		}
	})

	t.Run("remove participant", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		sf := NewMockSecretFriendFacade(ctrl)
		repo.EXPECT().GetParticipant(sfID, userID).Return(participantEntity, nil)
		repo.EXPECT().RemoveParticipant(sfID, userID).Return(nil)
		uc := New(associatedUser, repo, sf)
		if err := uc.RemoveParticipant(sfID); err != nil {
			t.Fatalf("RemoveParticipant() error = %v, want nil", err)
		}
	})

	t.Run("check participant exists proxies repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := NewMockRepository(ctrl)
		sf := NewMockSecretFriendFacade(ctrl)
		repo.EXPECT().GetParticipant(sfID, userID).Return(participantEntity, nil)
		uc := New(associatedUser, repo, sf)
		got, err := uc.CheckParticipantExists(sfID, userID)
		if err != nil || got.ID != participantID {
			t.Fatalf("CheckParticipantExists() = (%+v, %v), want participant and nil", got, err)
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
