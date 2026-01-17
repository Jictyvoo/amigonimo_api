package drawfriends

import (
	"fmt"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Repository interface {
	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateSecretFriend(sf *entities.SecretFriend) error
	GetParticipantsCount(secretFriendID entities.HexID) (int, error)
	GetDrawResultForUser(secretFriendID, userID entities.HexID) (entities.DrawResultItem, error)
	SaveDrawResults(secretFriendID entities.HexID, results []entities.DrawResultItem) error
}

type UseCase struct {
	repo     Repository
	drawServ *drawserv.Service
}

func New(repo Repository, drawServ *drawserv.Service) *UseCase {
	return &UseCase{
		repo:     repo,
		drawServ: drawServ,
	}
}

func (uc *UseCase) Execute(input ExecuteInput) (ExecuteOutput, error) {
	sf, err := uc.repo.GetSecretFriendByID(input.SecretFriendID)
	if err != nil {
		return ExecuteOutput{}, fmt.Errorf("get secret friend: %w", err)
	}

	if sf.Status == entities.StatusDrawn || sf.Status == entities.StatusClosed {
		return ExecuteOutput{}, fmt.Errorf("secret friend already drawn")
	}

	drawResult, err := uc.drawServ.ExecuteDraw(
		drawserv.DrawInput{
			Participants: sf.Participants,
		},
	)
	if err != nil {
		return ExecuteOutput{}, fmt.Errorf("execute drawfriends algorithm: %w", err)
	}

	if err := uc.repo.SaveDrawResults(sf.ID, drawResult.Pairs); err != nil {
		return ExecuteOutput{}, fmt.Errorf("save drawfriends results: %w", err)
	}

	sf.Status = entities.StatusDrawn
	sf.UpdatedAt = time.Now()

	if err := uc.repo.UpdateSecretFriend(&sf); err != nil {
		return ExecuteOutput{}, fmt.Errorf("update status after drawfriends: %w", err)
	}

	count, err := uc.repo.GetParticipantsCount(sf.ID)
	if err != nil {
		return ExecuteOutput{}, fmt.Errorf("get participants count: %w", err)
	}

	return ExecuteOutput{ParticipantCount: count}, nil
}

func (uc *UseCase) GetResult(input GetResultInput) (entities.DrawResultItem, error) {
	result, err := uc.repo.GetDrawResultForUser(input.SecretFriendID, input.UserID)
	if err != nil {
		return entities.DrawResultItem{}, fmt.Errorf("get personal drawfriends result: %w", err)
	}
	return result, nil
}
