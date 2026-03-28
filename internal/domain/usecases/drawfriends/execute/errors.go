package execute

import "errors"

var (
	ErrNoValidDraw         = errors.New("no valid drawfriends found after maximum attempts")
	ErrInsufficientPlayers = errors.New("at least 3 participants are required for a drawfriends")
)
