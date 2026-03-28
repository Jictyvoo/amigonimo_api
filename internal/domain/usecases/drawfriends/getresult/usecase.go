package getresult

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Input struct {
	SecretFriendID entities.HexID
}

type Output = entities.DrawResultItem

type UseCase struct {
	repo           Repository
	associatedUser entities.User
}

func New(associatedUser entities.User, repo Repository) UseCase {
	return UseCase{
		repo:           repo,
		associatedUser: associatedUser,
	}
}

func (uc UseCase) Execute(input Input) (Output, error) {
	result, err := uc.repo.GetDrawResultForUser(input.SecretFriendID, uc.associatedUser.ID)
	if err != nil {
		return Output{}, apperr.From(
			"draw_result_not_found",
			"draw result not found",
			err,
		)
	}
	return result, nil
}
