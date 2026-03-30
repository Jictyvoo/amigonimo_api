package fixturesets

import (
	"database/sql"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
	"github.com/jictyvoo/amigonimo_api/test/internal/fixtures"
)

// SecretFriendRoutes holds a 3-participant open event ready for route testing.
type SecretFriendRoutes struct {
	Owner           *User
	Target          *User
	Other           *User
	OpenEvent       *genmodels.SecretFriend
	OpenOwnerEntry  *genmodels.Participant
	OpenTargetEntry *genmodels.Participant
	OpenOtherEntry  *genmodels.Participant
}

func NewSecretFriendRoutes(owner, target, other *User) SecretFriendRoutes {
	openEvent := fixtures.NewSecretFriend().
		WithOwner(owner.User).
		WithName("Open Routes Event").
		Build()
	openEvent.InviteCode = "openrt01"
	openEvent.Location = sql.NullString{String: "Original Place", Valid: true}
	openEvent.Datetime = sql.NullTime{Time: time.Now().Add(48 * time.Hour), Valid: true}

	openOwnerEntry := fixtures.NewParticipant().
		WithUser(owner.User).
		WithSecretFriend(openEvent).
		Build()
	openTargetEntry := fixtures.NewParticipant().
		WithUser(target.User).
		WithSecretFriend(openEvent).
		Build()
	openOtherEntry := fixtures.NewParticipant().
		WithUser(other.User).
		WithSecretFriend(openEvent).
		Build()

	return SecretFriendRoutes{
		Owner:           owner,
		Target:          target,
		Other:           other,
		OpenEvent:       openEvent,
		OpenOwnerEntry:  openOwnerEntry,
		OpenTargetEntry: openTargetEntry,
		OpenOtherEntry:  openOtherEntry,
	}
}

func (s SecretFriendRoutes) Seedables() []any {
	return []any{
		s.Owner.User,
		s.Owner.Profile,
		s.Target.User,
		s.Target.Profile,
		s.Other.User,
		s.Other.Profile,
		s.OpenEvent,
		s.OpenOwnerEntry,
		s.OpenTargetEntry,
		s.OpenOtherEntry,
	}
}

func (s SecretFriendRoutes) OpenEventID() entities.HexID {
	return mustHexIDFromBytes(s.OpenEvent.ID)
}

func (s SecretFriendRoutes) OpenOwnerEntryID() entities.HexID {
	return mustHexIDFromBytes(s.OpenOwnerEntry.ID)
}

func (s SecretFriendRoutes) OpenTargetEntryID() entities.HexID {
	return mustHexIDFromBytes(s.OpenTargetEntry.ID)
}

func (s SecretFriendRoutes) OpenOtherEntryID() entities.HexID {
	return mustHexIDFromBytes(s.OpenOtherEntry.ID)
}
