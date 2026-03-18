package secretfriend

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

func (uc *UseCase) GetInviteInfo(code string) (entities.SecretFriend, error) {
	sf, err := uc.repo.GetSecretFriendByInviteCode(code)
	if err != nil {
		return entities.SecretFriend{}, apperr.From(
			"secret_friend_invite_not_found",
			"invite code not found",
			err,
		)
	}

	return sf, nil
}
