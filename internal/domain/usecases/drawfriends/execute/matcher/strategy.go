package matcher

import "errors"

// Priority ranks strategy results. Lower values are preferred.
type Priority uint8

const (
	PriorityChainClose Priority = iota
	PriorityGreedy
	PriorityReverseGreedy
	PriorityChainRestart
	PriorityRotation
)

// DrawStrategy defines the interface every draw algorithm must satisfy.
type DrawStrategy interface {
	// Execute attempts to find a valid set of pairings for the given participants.
	// The orchestrator already pre-computes each participant's AllowedReceivers list.
	Execute(participants []Participant) ([]Pairing, error)

	// ResultPriority returns the strategy's fixed priority ranking.
	ResultPriority() Priority
}

var (
	ErrNoValidDraw         = errors.New("no valid draw found after exhausting all strategies")
	ErrInsufficientPlayers = errors.New("at least 3 participants are required for a draw")
)
