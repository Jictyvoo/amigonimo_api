package matcher

// ChainCloseStrategy attempts to build a single Hamiltonian cycle through all
// participants. It starts from the most constrained participant and follows
// the chain, always moving to the assigned receiver as the next giver.
//
// This produces the ideal secret-friend draw: a single ring where every
// participant gives to one person and receives from another, with no sub-groups.
//
// Backtracking is stack-based (no recursion), limited to n*2 total backtracks.
type ChainCloseStrategy struct {
	baseStrategy
}

func (ChainCloseStrategy) ResultPriority() Priority { return PriorityChainClose }

func (cc ChainCloseStrategy) Execute(participants []Participant) ([]Pairing, error) {
	graph, n, err := cc.setup(participants, true)
	if err != nil {
		return nil, err
	}

	cs := newChainSearch(n, 0)

	for cs.active() {
		if cs.depth() == n {
			if tryCloseCycle(&graph, &cs) {
				return graph.chainToPairings(cs.chain), nil
			}
			continue
		}

		candidates := graph.unusedCandidates(cs.currentIdx(), cs.used)
		step := cs.depth() - 1
		nextIdx, ok := cs.pickCandidate(step, candidates)
		if ok {
			cs.advance(nextIdx)
		} else if !cs.retreat() {
			break
		}
	}

	return nil, ErrNoValidDraw
}

// tryCloseCycle checks whether the last participant in the chain can give to
// the first (closing the Hamiltonian cycle). If not, it retreats.
func tryCloseCycle(graph *constraintGraph, cs *chainSearch) bool {
	lastIdx := cs.currentIdx()
	startID := graph.sorted[cs.chain[0]].ID
	if graph.canReceive(lastIdx, startID) {
		return true
	}
	cs.retreat()
	return false
}
