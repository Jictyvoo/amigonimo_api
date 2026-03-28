package execute

import (
	"context"
	"errors"

	"github.com/jictyvoo/amigonimo_api/internal/domain/apperr"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type Input struct {
	SecretFriendID entities.HexID
}

type Output struct {
	ParticipantCount int
}

type UseCase struct {
	repo               Repository
	secretFriendFacade SecretFriendFacade
	friendMatcher      DrawFriendMatcher
}

func New(repo Repository, sfFacade SecretFriendFacade) UseCase {
	return UseCase{
		repo:               repo,
		secretFriendFacade: sfFacade,
		friendMatcher:      NewDrawMatcher(),
	}
}

func (uc UseCase) Execute(input Input) (output Output, err error) {
	var sf entities.SecretFriend
	if sf, err = uc.secretFriendFacade.GetSecretFriendByID(input.SecretFriendID); err != nil {
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

	var drawResult DrawOutput
	if drawResult, err = uc.friendMatcher.ExecuteDraw(
		DrawInput{Participants: sf.Participants},
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

	if err = uc.secretFriendFacade.UpdateStatus(sf.ID, entities.StatusDrawn); err != nil {
		return Output{}, apperr.From(
			"secret_friend_status_update_failed",
			"failed to update secret friend status",
			err,
		)
	}

	return Output{ParticipantCount: len(sf.Participants)}, nil
}
