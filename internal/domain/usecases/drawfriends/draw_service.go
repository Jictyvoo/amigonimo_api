package drawfriends

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/execute"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/getresult"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type (
	Repository interface {
		getresult.Repository
		execute.Repository
	}
	Facades interface {
		execute.SecretFriendFacade
	}
)

type Service struct {
	executeUseCase   execute.UseCase
	getResultUseCase getresult.UseCase
}

func New(
	associatedUser entities.User, repo Repository, facades Facades,
) Service {
	return Service{
		executeUseCase:   execute.New(repo, facades, drawserv.New()),
		getResultUseCase: getresult.New(associatedUser, repo),
	}
}

func (uc Service) Execute(secretFriendID entities.HexID) (int, error) {
	result, err := uc.executeUseCase.Execute(execute.Input{SecretFriendID: secretFriendID})
	if err != nil {
		return 0, err
	}

	return result.ParticipantCount, err
}

func (uc Service) GetResult(secretFriendID entities.HexID) (getresult.Output, error) {
	return uc.getResultUseCase.Execute(getresult.Input{SecretFriendID: secretFriendID})
}
