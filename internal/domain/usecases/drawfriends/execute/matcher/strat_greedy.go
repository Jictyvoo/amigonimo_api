package matcher

// GreedyStrategy processes givers from most-constrained to least-constrained,
// always picking the first available receiver. It is the simplest backtracking
// approach and serves as a baseline for the other strategies.
type GreedyStrategy struct {
	baseStrategy
}

func (GreedyStrategy) ResultPriority() Priority { return PriorityGreedy }

func (g GreedyStrategy) Execute(participants []Participant) ([]Pairing, error) {
	graph, n, err := g.setup(participants, true)
	if err != nil {
		return nil, err
	}

	as := newAssignmentSearch(n)
	for as.active(n) {
		candidates := graph.unusedCandidates(as.step, as.usedAsReceiver)
		recvIdx, ok := as.pickCandidate(as.step, candidates)
		if ok {
			as.assign(recvIdx)
		} else if !as.retreat() {
			break
		}
	}

	if !as.isComplete(n) {
		return nil, ErrNoValidDraw
	}
	return graph.assignmentToPairings(as.assignment), nil
}
