package drawfriends

import (
	"context"
	"errors"
	"fmt"

	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

type Repository interface {
	dbrock.Transactioner

	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateSecretFriend(sf *entities.SecretFriend) error
	ListParticipants(
		ctx context.Context,
		secretFriendID entities.HexID,
	) ([]entities.Participant, error)
	GetDrawResultForUser(secretFriendID, userID entities.HexID) (entities.DrawResultItem, error)
	SaveDrawResults(secretFriendID entities.HexID, results []entities.DrawResultItem) error
}

type UseCase struct {
	repo     Repository
	drawServ drawserv.Service
}

func New(repo Repository, drawServ drawserv.Service) UseCase {
	return UseCase{
		repo:     repo,
		drawServ: drawServ,
	}
}

func (uc *UseCase) Execute(input ExecuteInput) (output ExecuteOutput, err error) {
	var sf entities.SecretFriend
	if sf, err = uc.repo.GetSecretFriendByID(input.SecretFriendID); err != nil {
		return ExecuteOutput{}, fmt.Errorf("get secret friend: %w", err)
	}

	if sf.Status == entities.StatusDrawn || sf.Status == entities.StatusClosed {
		return ExecuteOutput{}, fmt.Errorf("secret friend already drawn")
	}

	var drawResult drawserv.DrawOutput
	if drawResult, err = uc.drawServ.ExecuteDraw(
		drawserv.DrawInput{
			Participants: sf.Participants,
		},
	); err != nil {
		return ExecuteOutput{}, fmt.Errorf("execute draw algorithm: %w", err)
	}

	onFinishTx, txErr := uc.repo.BeginTx(context.Background(), nil)
	if txErr != nil {
		return ExecuteOutput{}, fmt.Errorf("failed to begin tx: %w", txErr)
	}

	defer func() { // Finish transaction
		txErr = onFinishTx(err == nil)
		if txErr != nil {
			err = errors.Join(err, txErr)
		}
	}()

	if err = uc.repo.SaveDrawResults(sf.ID, drawResult.Pairs); err != nil {
		return ExecuteOutput{}, fmt.Errorf("save draw results: %w", err)
	}

	sf.Status = entities.StatusDrawn
	if err = uc.repo.UpdateSecretFriend(&sf); err != nil {
		return ExecuteOutput{}, fmt.Errorf("update status after draw: %w", err)
	}

	return ExecuteOutput{ParticipantCount: len(sf.Participants)}, nil
}

func (uc *UseCase) GetResult(input GetResultInput) (entities.DrawResultItem, error) {
	result, err := uc.repo.GetDrawResultForUser(input.SecretFriendID, input.UserID)
	if err != nil {
		return entities.DrawResultItem{}, fmt.Errorf("get personal draw result: %w", err)
	}
	return result, nil
}
