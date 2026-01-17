package drawfriends

import "github.com/jictyvoo/amigonimo_api/internal/entities"

type ExecuteInput struct {
	SecretFriendID entities.HexID
}

type ExecuteOutput struct {
	ParticipantCount int
}

type GetResultInput struct {
	SecretFriendID entities.HexID
	UserID         entities.HexID
}
