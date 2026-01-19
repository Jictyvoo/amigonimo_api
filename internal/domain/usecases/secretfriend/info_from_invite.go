package secretfriend

import (
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) GetInviteInfo(code string) (entities.SecretFriend, error) {
	return uc.repo.GetSecretFriendByInviteCode(code)
}
