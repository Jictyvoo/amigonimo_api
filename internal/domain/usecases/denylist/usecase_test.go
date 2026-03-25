package denylist

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

type secretFriendFacadeAdapter struct {
	ports.Facade
	*MocksecretFriendFacadePort
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

func TestDenylistUseCase(t *testing.T) {
	sfID := mustHexID(t)
	userID := mustHexID(t)
	targetID := mustHexID(t)
	participantID := mustHexID(t)
	associatedUser := entities.User{ID: userID}
	ref := ParticipantRef{ParticipantID: participantID, UserID: userID, SecretFriendID: sfID}

	t.Run(
		"get denylist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			provider := NewFacadeProvider(
				participantFacadeAdapter{
					MockparticipantFacadePort: NewMockparticipantFacadePort(ctrl),
				},
				secretFriendFacadeAdapter{
					MocksecretFriendFacadePort: NewMocksecretFriendFacadePort(ctrl),
				},
			)
			repo.EXPECT().
				GetDenyListByParticipant(ParticipantRef{UserID: userID, SecretFriendID: sfID}).
				Return([]entities.DeniedUser{{ID: mustHexID(t)}}, nil)
			uc := New(associatedUser, repo, provider)
			items, err := uc.GetDenyList(sfID)
			if err != nil || len(items) != 1 {
				t.Fatalf("GetDenyList() = (%v, %v), want one item and nil", items, err)
			}
		},
	)

	t.Run(
		"add entry validations and success", func(t *testing.T) {
			tests := []struct {
				name    string
				setup   func(*MockRepository, *MockparticipantFacadePort, *MocksecretFriendFacadePort)
				target  entities.HexID
				wantErr error
			}{
				{
					name:    "self deny is invalid",
					target:  userID,
					wantErr: errors.New("self"),
				},
				{
					name:   "capacity reached",
					target: targetID,
					setup: func(repo *MockRepository, participant *MockparticipantFacadePort, sf *MocksecretFriendFacadePort) {
						participant.EXPECT().
							CheckParticipantInSecretFriend(sfID, userID).
							Return(entities.Participant{ID: participantID}, nil)
						participant.EXPECT().
							CheckParticipantInSecretFriend(sfID, targetID).
							Return(entities.Participant{ID: mustHexID(t)}, nil)
						repo.EXPECT().
							GetDenyListByParticipant(ref).
							Return([]entities.DeniedUser{{ID: mustHexID(t)}, {ID: mustHexID(t)}}, nil)
						sf.EXPECT().
							GetSecretFriendByID(sfID).
							Return(entities.SecretFriend{ID: sfID, MaxDenyListSize: 2}, nil)
					},
					wantErr: errors.New("capacity"),
				},
				{
					name:   "success",
					target: targetID,
					setup: func(repo *MockRepository, participant *MockparticipantFacadePort, sf *MocksecretFriendFacadePort) {
						participant.EXPECT().
							CheckParticipantInSecretFriend(sfID, userID).
							Return(entities.Participant{ID: participantID}, nil)
						participant.EXPECT().
							CheckParticipantInSecretFriend(sfID, targetID).
							Return(entities.Participant{ID: mustHexID(t)}, nil)
						repo.EXPECT().GetDenyListByParticipant(ref).Return(nil, nil)
						sf.EXPECT().
							GetSecretFriendByID(sfID).
							Return(entities.SecretFriend{ID: sfID, MaxDenyListSize: 2}, nil)
						repo.EXPECT().
							AddDenyListEntry(ref, targetID).
							Return(entities.DeniedUser{ID: mustHexID(t)}, nil)
					},
				},
			}
			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						ctrl := gomock.NewController(t)
						repo := NewMockRepository(ctrl)
						participant := NewMockparticipantFacadePort(ctrl)
						sf := NewMocksecretFriendFacadePort(ctrl)
						if tt.setup != nil {
							tt.setup(repo, participant, sf)
						}
						uc := New(
							associatedUser,
							repo,
							NewFacadeProvider(
								participantFacadeAdapter{MockparticipantFacadePort: participant},
								secretFriendFacadeAdapter{MocksecretFriendFacadePort: sf},
							),
						)
						_, err := uc.AddEntry(sfID, tt.target)
						if tt.wantErr != nil {
							wantCode := "denylist_self_entry"
							if tt.name == "capacity reached" {
								wantCode = "denylist_capacity_reached"
							}
							requireAppErrorCode(t, err, wantCode)
							return
						}
						if err != nil {
							t.Fatalf("AddEntry() error = %v, want nil", err)
						}
					},
				)
			}
		},
	)

	t.Run(
		"remove entry", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := NewMockRepository(ctrl)
			participant := NewMockparticipantFacadePort(ctrl)
			sf := NewMocksecretFriendFacadePort(ctrl)
			participant.EXPECT().
				CheckParticipantInSecretFriend(sfID, userID).
				Return(entities.Participant{ID: participantID}, nil)
			repo.EXPECT().RemoveDenyListEntry(ref, targetID).Return(nil)
			uc := New(
				associatedUser,
				repo,
				NewFacadeProvider(
					participantFacadeAdapter{MockparticipantFacadePort: participant},
					secretFriendFacadeAdapter{MocksecretFriendFacadePort: sf},
				),
			)
			if err := uc.RemoveEntry(sfID, targetID); err != nil {
				t.Fatalf("RemoveEntry() error = %v, want nil", err)
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
