package execute

import (
	"context"
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/domain/services/drawserv"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/pkg/dbrock"
)

type Repository interface {
	dbrock.Transactioner

	GetSecretFriendByID(id entities.HexID) (entities.SecretFriend, error)
	UpdateSecretFriend(sf *entities.SecretFriend) error
	SaveDrawResults(secretFriendID entities.HexID, results []entities.DrawResultItem) error
}

type Input struct {
	SecretFriendID entities.HexID
}

type Output struct {
	ParticipantCount int
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

func (uc UseCase) Execute(input Input) (output Output, err error) {
	var sf entities.SecretFriend
	if sf, err = uc.repo.GetSecretFriendByID(input.SecretFriendID); err != nil {
		return Output{}, apperr.From(
			"secret_friend_not_found",
			"secret friend not found",
			err,
		)
	}

	if sf.Status == entities.StatusDrawn || sf.Status == entities.StatusClosed {
		return Output{}, apperr.Conflict(
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
		return Output{}, apperr.InternalError(
			"draw_execution_failed",
			"failed to execute draw",
			err,
		)
	}

	onFinishTx, txErr := uc.repo.BeginTx(context.Background(), nil)
	if txErr != nil {
		return Output{}, apperr.InternalError(
			"draw_transaction_start_failed",
			"failed to start draw transaction",
			txErr,
		)
	}

	defer func() {
		txErr = onFinishTx(err == nil)
		if txErr != nil {
			err = errors.Join(err, txErr)
		}
	}()

	if err = uc.repo.SaveDrawResults(sf.ID, drawResult.Pairs); err != nil {
		return Output{}, apperr.From(
			"draw_result_save_failed",
			"failed to save draw results",
			err,
		)
	}

	sf.Status = entities.StatusDrawn
	if err = uc.repo.UpdateSecretFriend(&sf); err != nil {
		return Output{}, apperr.From(
			"secret_friend_status_update_failed",
			"failed to update secret friend status",
			err,
		)
	}

	return Output{ParticipantCount: len(sf.Participants)}, nil
}
