package matcher

import "slices"

// ReverseGreedyStrategy is the dual of GreedyStrategy. It processes givers
// from least-constrained to most-constrained. For each giver, it picks the
// receiver that has the fewest remaining allowed givers — the receiver most
// at risk of becoming unassignable.
//
// This explores a fundamentally different region of the search space than
// the standard greedy approach.
type ReverseGreedyStrategy struct {
	baseStrategy
}

func (ReverseGreedyStrategy) ResultPriority() Priority { return PriorityReverseGreedy }

func (rg ReverseGreedyStrategy) Execute(participants []Participant) ([]Pairing, error) {
	graph, n, err := rg.setup(participants, false)
	if err != nil {
		return nil, err
	}
	as := newAssignmentSearch(n)
	inboundCount := buildInboundCounts(&graph)

	for as.active(n) {
		candidates := candidatesByInboundRisk(&graph, as.step, as.usedAsReceiver, inboundCount)
		recvIdx, ok := as.pickCandidate(as.step, candidates)
		if ok {
			adjustInboundCounts(&graph, as.step, recvIdx, inboundCount, -1)
			as.assign(recvIdx)
		} else if !reverseGreedyRetreat(&as, &graph, inboundCount) {
			break
		}
	}

	if !as.isComplete(n) {
		return nil, ErrNoValidDraw
	}
	return graph.assignmentToPairings(as.assignment), nil
}

// buildInboundCounts computes how many givers can choose each participant as receiver.
func buildInboundCounts(graph *constraintGraph) []int {
	n := graph.size()
	inboundCount := make([]int, n)
	for gi := range n {
		for _, rid := range graph.sorted[gi].AllowedReceivers {
			if ri, ok := graph.idxByID[rid]; ok {
				inboundCount[ri]++
			}
		}
	}
	return inboundCount
}

// adjustInboundCounts modifies inbound counts for a receiver from givers after the
// current step. delta is +1 (restore) or -1 (consume).
func adjustInboundCounts(
	graph *constraintGraph,
	fromStep int,
	recvIdx int,
	inboundCount []int,
	delta int,
) {
	recvID := graph.sorted[recvIdx].ID
	for gi := fromStep + 1; gi < graph.size(); gi++ {
		if graph.canReceive(gi, recvID) {
			inboundCount[recvIdx] += delta
		}
	}
}

// reverseGreedyRetreat undoes the current and previous assignments, restoring
// inbound counts for both. Returns false if the search is exhausted.
func reverseGreedyRetreat(as *assignmentSearch, graph *constraintGraph, inboundCount []int) bool {
	// Restore for current step's assignment (if any).
	if as.assignment[as.step] >= 0 {
		adjustInboundCounts(graph, as.step, as.assignment[as.step], inboundCount, +1)
	}
	prevStep := as.step - 1
	if prevStep >= 0 {
		adjustInboundCounts(graph, prevStep, as.assignment[prevStep], inboundCount, +1)
	}
	return as.retreat()
}

// candidatesByInboundRisk returns candidates sorted by ascending inbound count
// (most at-risk receivers first), then by ID for determinism.
func candidatesByInboundRisk(
	graph *constraintGraph,
	giverStep int,
	usedAsReceiver []bool,
	inboundCount []int,
) []int {
	candidates := graph.unusedCandidates(giverStep, usedAsReceiver)
	slices.SortFunc(
		candidates, func(a, b int) int {
			diff := inboundCount[a] - inboundCount[b]
			if diff != 0 {
				return diff
			}
			return graph.sorted[a].ID.Compare(graph.sorted[b].ID)
		},
	)
	return candidates
}
