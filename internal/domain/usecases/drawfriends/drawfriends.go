package drawfriends

import (
	"context"
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

//go:generate go tool -modfile=../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=drawfriends github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends Repository

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
		return ExecuteOutput{}, apperr.From(
			"secret_friend_not_found",
			"secret friend not found",
			err,
		)
	}

	if sf.Status == entities.StatusDrawn || sf.Status == entities.StatusClosed {
		return ExecuteOutput{}, apperr.Conflict(
			"secret_friend_already_drawn",
			"secret friend draw has already been completed",
			nil,
		)
	}

	var drawResult drawserv.DrawOutput
	if drawResult, err = uc.drawServ.ExecuteDraw(
		drawserv.DrawInput{
			Participants: sf.Participants,
		},
	); err != nil {
		return ExecuteOutput{}, apperr.InternalError(
			"draw_execution_failed",
			"failed to execute draw",
			err,
		)
	}

	onFinishTx, txErr := uc.repo.BeginTx(context.Background(), nil)
	if txErr != nil {
		return ExecuteOutput{}, apperr.InternalError(
			"draw_transaction_start_failed",
			"failed to start draw transaction",
			txErr,
		)
	}

	defer func() { // Finish transaction
		txErr = onFinishTx(err == nil)
		if txErr != nil {
			err = errors.Join(err, txErr)
		}
	}()

	if err = uc.repo.SaveDrawResults(sf.ID, drawResult.Pairs); err != nil {
		return ExecuteOutput{}, apperr.From(
			"draw_result_save_failed",
			"failed to save draw results",
			err,
		)
	}

	sf.Status = entities.StatusDrawn
	if err = uc.repo.UpdateSecretFriend(&sf); err != nil {
		return ExecuteOutput{}, apperr.From(
			"secret_friend_status_update_failed",
			"failed to update secret friend status",
			err,
		)
	}

	return ExecuteOutput{ParticipantCount: len(sf.Participants)}, nil
}

func (uc *UseCase) GetResult(input GetResultInput) (entities.DrawResultItem, error) {
	result, err := uc.repo.GetDrawResultForUser(input.SecretFriendID, input.UserID)
	if err != nil {
		return entities.DrawResultItem{}, apperr.From(
			"draw_result_not_found",
			"draw result not found",
			err,
		)
	}
	return result, nil
}
