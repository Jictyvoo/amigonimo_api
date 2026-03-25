package verify

import "github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"

type UseCase struct {
	userRepository Repository
}

func New(userRepository Repository) UseCase {
	return UseCase{userRepository: userRepository}
}

func (uc UseCase) Execute(code string) error {
	user, err := uc.userRepository.GetUserByVerificationCode(code)
	if err != nil || user.ID.IsEmpty() {
		return autherrs.ErrVerificationCode
	}
	if err = uc.userRepository.SetUserVerified(user.ID); err != nil {
		return autherrs.NewErrSetVerification(err)
	}
	return nil
}
