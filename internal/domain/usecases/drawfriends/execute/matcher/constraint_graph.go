package matcher

// constraintGraph holds the pre-computed lookup structures that every strategy needs.
// Build it once and pass to the strategy's core loop.
type constraintGraph struct {
	sorted      []Participant
	allowedSets []map[ParticipantID]struct{}
	idxByID     map[ParticipantID]int
}

func newConstraintGraph(participants []Participant, ascending bool) constraintGraph {
	sorted := sortByConstraint(participants, ascending)
	n := len(sorted)
	g := constraintGraph{
		sorted:      sorted,
		allowedSets: make([]map[ParticipantID]struct{}, n),
		idxByID:     make(map[ParticipantID]int, n),
	}
	for i, p := range sorted {
		g.allowedSets[i] = buildAllowedSet(p)
		g.idxByID[p.ID] = i
	}
	return g
}

func (g *constraintGraph) size() int { return len(g.sorted) }

// canReceive checks whether the giver at giverIdx is allowed to give to receiverID.
func (g *constraintGraph) canReceive(giverIdx int, receiverID ParticipantID) bool {
	_, ok := g.allowedSets[giverIdx][receiverID]
	return ok
}

// unusedCandidates returns indices of participants that are in the giver's allowed
// set and not yet marked as used. Results are in AllowedReceivers order (which is
// deterministic since participants are pre-sorted).
func (g *constraintGraph) unusedCandidates(participantIdx int, used []bool) []int {
	p := g.sorted[participantIdx]
	candidates := make([]int, 0, len(p.AllowedReceivers))
	for _, rid := range p.AllowedReceivers {
		idx := g.idxByID[rid]
		if !used[idx] {
			candidates = append(candidates, idx)
		}
	}
	return candidates
}

// assignmentToPairings converts a step→receiver-index mapping into Pairing structs.
func (g *constraintGraph) assignmentToPairings(assignment []int) []Pairing {
	n := len(assignment)
	pairs := make([]Pairing, n)
	for i := range n {
		pairs[i] = Pairing{
			GiverID:    g.sorted[i].ID,
			ReceiverID: g.sorted[assignment[i]].ID,
		}
	}
	return pairs
}
