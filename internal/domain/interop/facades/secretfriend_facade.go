package facades

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/interop/ports"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/denylist"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/execute"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/participant"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/secretfriend"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// Ensure SecretFriendFacade implements required interfaces.
var (
	_ participant.SecretFriendFacade = (*SecretFriendFacade)(nil)
	_ denylist.SecretFriendFacade    = (*SecretFriendFacade)(nil)
	_ execute.SecretFriendFacade     = (*SecretFriendFacade)(nil)
)

type SecretFriendFacade struct {
	ports.Facade
	uc secretfriend.UseCase
}

func NewSecretFriendFacade(uc secretfriend.UseCase) *SecretFriendFacade {
	return &SecretFriendFacade{uc: uc}
}

func (f *SecretFriendFacade) GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error) {
	return f.uc.Get(id)
}

func (f *SecretFriendFacade) CheckUserIsOwner(sfID entities.HexID) (bool, error) {
	return f.uc.CheckUserIsOwner(sfID)
}

func (f *SecretFriendFacade) UpdateStatus(
	id entities.HexID, status entities.SecretFriendStatus,
) error {
	return f.uc.Update(secretfriend.UpdateInput{ID: id, Status: status})
}
