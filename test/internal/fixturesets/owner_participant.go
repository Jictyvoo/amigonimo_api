package fixturesets

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

type OwnerParticipant struct {
	Owner        *User
	Participant  *User
	SecretFriend *genmodels.SecretFriend
	OwnerEntry   *genmodels.Participant
	UserEntry    *genmodels.Participant
}

func NewOwnerParticipant(owner, participant *User, eventName string) *OwnerParticipant {
	event := fixtures.NewSecretFriend().
		WithOwner(owner.User).
		WithName(eventName).
		Build()
	ownerEntry := fixtures.NewParticipant().
		WithUser(owner.User).
		WithSecretFriend(event).
		Build()
	participantEntry := fixtures.NewParticipant().
		WithUser(participant.User).
		WithSecretFriend(event).
		Build()

	return &OwnerParticipant{
		Owner:        owner,
		Participant:  participant,
		SecretFriend: event,
		OwnerEntry:   ownerEntry,
		UserEntry:    participantEntry,
	}
}

// WithMaxDenyListSize sets the MaxDenyListSize on the seeded event.
func (s *OwnerParticipant) WithMaxDenyListSize(n uint8) *OwnerParticipant {
	s.SecretFriend.MaxDenyListSize = n
	return s
}

func (s *OwnerParticipant) Seedables() []any {
	return []any{
		s.Owner.User,
		s.Owner.Profile,
		s.Participant.User,
		s.Participant.Profile,
		s.SecretFriend,
		s.OwnerEntry,
		s.UserEntry,
	}
}

func (s *OwnerParticipant) SecretFriendID() entities.HexID {
	return mustHexIDFromBytes(s.SecretFriend.ID)
}
