package checkrecovery

import (
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

const resetCodeExpiration = 30 * time.Minute

type UseCase struct {
	userRepository Repository
}

func New(userRepository Repository) UseCase {
	return UseCase{userRepository: userRepository}
}

func (uc UseCase) Execute(
	identifier, recoveryCode string,
) (userID entities.HexID, err error) {
	user, err := uc.userRepository.GetUserByRecovery(
		identifier, recoveryCode, time.Now().Add(-resetCodeExpiration),
	)
	if err != nil || user.ID.IsEmpty() {
		return entities.HexID{}, autherrs.ErrUserRecoveryNotFound
	}

	return user.ID, nil
}
