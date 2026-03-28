package matcher

import (
	"context"
	"sync"
)

// Orchestrator runs all draw strategies concurrently and picks the best result.
// It satisfies the DrawStrategy interface itself so callers see a single entry point.
type Orchestrator struct {
	strategies []DrawStrategy
}

// NewOrchestrator creates an Orchestrator pre-loaded with all five strategies.
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		strategies: []DrawStrategy{
			GreedyStrategy{},
		},
	}
}

type strategyResult struct {
	pairs    []Pairing
	priority Priority
}

// Execute builds the allowed graph once, then fans out to all strategies.
// The result with the best (lowest) priority wins. If the best possible
// priority (ChainClose) arrives, remaining strategies are cancelled early.
func (o *Orchestrator) Execute(participants []Participant) ([]Pairing, error) {
	n := len(participants)
	if n < 3 {
		return nil, ErrInsufficientPlayers
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resultCh := make(chan strategyResult, len(o.strategies))
	var wg sync.WaitGroup

	for _, s := range o.strategies {
		wg.Add(1)
		go func(strategy DrawStrategy) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			pairs, err := strategy.Execute(participants)
			if err != nil {
				return
			}

			resultCh <- strategyResult{
				pairs:    pairs,
				priority: strategy.ResultPriority(),
			}

			// Early exit: best possible priority achieved.
			if strategy.ResultPriority() == PriorityChainClose {
				cancel()
			}
		}(s)
	}

	// Close channel once all goroutines finish.
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var best *strategyResult
	for r := range resultCh {
		if best == nil || r.priority < best.priority {
			best = &strategyResult{pairs: r.pairs, priority: r.priority}
		}
	}

	if best == nil {
		return nil, ErrNoValidDraw
	}
	return best.pairs, nil
}
