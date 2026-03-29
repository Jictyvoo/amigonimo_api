package matcher

import "slices"

// GreedyStrategy assigns receivers using a most-constrained-first approach.
// Participants are sorted by ascending number of allowed receivers. For each
// giver, the first available (not-yet-assigned-as-receiver) candidate is chosen.
//
// If stuck, it backtracks iteratively (stack-based, no recursion) up to n*2 times.
// Produces a valid derangement, but not necessarily a single cycle.
type GreedyStrategy struct{}

func (GreedyStrategy) ResultPriority() Priority { return PriorityGreedy }

func (GreedyStrategy) Execute(participants []Participant) ([]Pairing, error) {
	n := len(participants)
	if n < 3 {
		return nil, ErrInsufficientPlayers
	}

	graph := newConstraintGraph(participants, true)
	as := newAssignmentSearch(n)

	for as.active(n) {
		candidates := sortedCandidatesByConstraint(&graph, as.step, as.usedAsReceiver)
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

// sortedCandidatesByConstraint returns indices of unused-as-receiver participants
// in the giver's allowed set, sorted by ascending constraint count for determinism.
func sortedCandidatesByConstraint(
	graph *constraintGraph,
	giverStep int,
	usedAsReceiver []bool,
) []int {
	candidates := graph.unusedCandidates(giverStep, usedAsReceiver)
	slices.SortFunc(candidates, func(a, b int) int {
		diff := len(graph.sorted[a].AllowedReceivers) - len(graph.sorted[b].AllowedReceivers)
		if diff != 0 {
			return diff
		}
		return graph.sorted[a].ID.Compare(graph.sorted[b].ID)
	})
	return candidates
}
