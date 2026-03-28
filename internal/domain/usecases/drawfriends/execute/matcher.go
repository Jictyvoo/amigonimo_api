package execute

import (
	"math/rand/v2"
	"time"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

type DrawFriendMatcher struct{}

func NewDrawMatcher() DrawFriendMatcher {
	return DrawFriendMatcher{}
}

type (
	DrawInput struct {
		Participants []entities.Participant
	}
	DrawOutput struct {
		Pairs []entities.DrawResultItem
	}
)

// ExecuteDraw generates a valid draw based on participants and their denylist.
// It uses a simple shuffle and check approach with a limited number of attempts.
func (s *DrawFriendMatcher) ExecuteDraw(input DrawInput) (DrawOutput, error) {
	if len(input.Participants) < 3 {
		return DrawOutput{}, ErrInsufficientPlayers
	}

	participants := input.Participants
	n := len(participants)

	// Pre-map denylist for faster lookup
	denyMaps := make(map[entities.HexID]map[entities.HexID]bool)
	for _, p := range participants {
		denyMaps[p.ID] = make(map[entities.HexID]bool)
		for _, d := range p.DenyList {
			denyMaps[p.ID][d.ID] = true
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	maxAttempts := 100
	for range maxAttempts {
		// Create a shuffled list of indices for targets
		targets := make([]int, n)
		for i := range targets {
			targets[i] = i
		}
		r.Shuffle(
			n, func(i, j int) {
				targets[i], targets[j] = targets[j], targets[i]
			},
		)

		valid := true
		pairs := make([]entities.DrawResultItem, n)

		for i := range n {
			giver := participants[i]
			targetIdx := targets[i]
			receiver := participants[targetIdx]

			// Constraint check:
			// 1. Cannot draw self
			// 2. Cannot draw someone in denylist
			if giver.ID == receiver.ID || denyMaps[giver.ID][receiver.ID] {
				valid = false
				break
			}

			pairs[i] = entities.DrawResultItem{
				Giver:    giver,
				Receiver: receiver,
			}
		}

		if valid {
			return DrawOutput{Pairs: pairs}, nil
		}
	}

	return DrawOutput{}, ErrNoValidDraw
}
