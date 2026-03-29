package matcher

// searchState is the common backtracking kernel shared by all search-based strategies.
// It tracks which candidate index has been tried at each step and how many
// backtracks have been consumed.
type searchState struct {
	candidatePos  []int
	backtracks    int
	maxBacktracks int
}

// pickCandidate selects the next untried candidate at the given step.
// It advances the cursor past the picked candidate so subsequent calls skip it.
func (s *searchState) pickCandidate(step int, candidates []int) (idx int, ok bool) {
	pos := s.candidatePos[step]
	if pos >= len(candidates) {
		return 0, false
	}
	s.candidatePos[step] = pos + 1
	return candidates[pos], true
}

// resetCandidates zeroes the candidate cursor at the given step.
func (s *searchState) resetCandidates(step int) {
	s.candidatePos[step] = 0
}

// recordBacktrack increments the counter and returns true if the budget is not exhausted.
func (s *searchState) recordBacktrack() bool {
	s.backtracks++
	return s.backtracks <= s.maxBacktracks
}

// assignmentSearch manages state for assignment-based strategies (Greedy, ReverseGreedy).
// Each step assigns a receiver index to a giver at that step position.
type assignmentSearch struct {
	searchState
	assignment     []int
	usedAsReceiver []bool
	step           int
}

func newAssignmentSearch(n int) assignmentSearch {
	assignment := make([]int, n)
	for i := range assignment {
		assignment[i] = -1
	}
	return assignmentSearch{
		searchState: searchState{
			candidatePos:  make([]int, n),
			maxBacktracks: n * 2,
		},
		assignment:     assignment,
		usedAsReceiver: make([]bool, n),
	}
}

// assign records the receiver for the current step and advances.
func (as *assignmentSearch) assign(recvIdx int) {
	as.assignment[as.step] = recvIdx
	as.usedAsReceiver[recvIdx] = true
	as.step++
}

// retreat undoes the current step and the previous one (the previous step needs
// to try its next candidate). Returns false if we've backed up past the start
// or exhausted the backtrack budget.
func (as *assignmentSearch) retreat() bool {
	as.resetCandidates(as.step)
	if as.assignment[as.step] >= 0 {
		as.usedAsReceiver[as.assignment[as.step]] = false
		as.assignment[as.step] = -1
	}
	as.step--
	if as.step < 0 {
		return false
	}
	as.usedAsReceiver[as.assignment[as.step]] = false
	as.assignment[as.step] = -1
	return as.recordBacktrack()
}

func (as *assignmentSearch) isComplete(n int) bool { return as.step == n }
func (as *assignmentSearch) active(n int) bool     { return as.step >= 0 && as.step < n }

// chainSearch manages state for chain-based strategies (ChainClose).
// It tracks a growing path of participant indices and which are already in the chain.
type chainSearch struct {
	searchState
	chain  []int
	used   []bool
	depth_ int
}

func newChainSearch(n, startIdx int) chainSearch {
	chain := make([]int, n)
	used := make([]bool, n)
	chain[0] = startIdx
	used[startIdx] = true
	return chainSearch{
		searchState: searchState{
			candidatePos:  make([]int, n),
			maxBacktracks: n * 2,
		},
		chain:  chain,
		used:   used,
		depth_: 1,
	}
}

func (cs *chainSearch) depth() int      { return cs.depth_ }
func (cs *chainSearch) active() bool    { return cs.depth_ > 0 }
func (cs *chainSearch) currentIdx() int { return cs.chain[cs.depth_-1] }

func (cs *chainSearch) advance(nextIdx int) {
	cs.chain[cs.depth_] = nextIdx
	cs.used[nextIdx] = true
	cs.depth_++
}

// retreat removes the last participant from the chain and resets its candidate cursor.
// Returns false when the search is exhausted (depth reaches zero or budget exceeded).
func (cs *chainSearch) retreat() bool {
	step := cs.depth_ - 1
	cs.resetCandidates(step)
	cs.used[cs.chain[step]] = false
	cs.depth_--
	if cs.depth_ <= 0 {
		return false
	}
	if !cs.recordBacktrack() {
		cs.depth_ = 0 // force active() = false so the loop exits
		return false
	}
	return true
}
