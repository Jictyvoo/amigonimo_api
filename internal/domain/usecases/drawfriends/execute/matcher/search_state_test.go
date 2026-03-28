package matcher

import "testing"

func TestSearchStatePickCandidate(t *testing.T) {
	tests := []struct {
		name       string
		startPos   int
		candidates []int
		wantIdx    int
		wantOk     bool
	}{
		{
			name:       "picks first candidate",
			candidates: []int{5, 10, 15},
			wantIdx:    5,
			wantOk:     true,
		},
		{
			name:       "picks from offset",
			startPos:   1,
			candidates: []int{5, 10, 15},
			wantIdx:    10,
			wantOk:     true,
		},
		{
			name:       "no candidates left",
			startPos:   3,
			candidates: []int{5, 10, 15},
			wantOk:     false,
		},
		{
			name:       "empty candidates",
			candidates: []int{},
			wantOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := searchState{candidatePos: []int{tt.startPos}}
			idx, ok := s.pickCandidate(0, tt.candidates)
			if ok != tt.wantOk {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOk)
			}
			if ok && idx != tt.wantIdx {
				t.Fatalf("idx = %d, want %d", idx, tt.wantIdx)
			}
		})
	}
}

func TestSearchStatePickCandidateAdvancesCursor(t *testing.T) {
	s := searchState{candidatePos: []int{0}}
	candidates := []int{5, 10, 15}

	idx1, _ := s.pickCandidate(0, candidates)
	idx2, _ := s.pickCandidate(0, candidates)

	if idx1 == idx2 {
		t.Fatal("consecutive picks should return different candidates")
	}
	if idx1 != 5 || idx2 != 10 {
		t.Fatalf("got %d, %d, want 5, 10", idx1, idx2)
	}
}

func TestSearchStateRecordBacktrack(t *testing.T) {
	s := searchState{maxBacktracks: 2}

	if !s.recordBacktrack() {
		t.Fatal("first backtrack should be within budget")
	}
	if !s.recordBacktrack() {
		t.Fatal("second backtrack should be within budget")
	}
	if s.recordBacktrack() {
		t.Fatal("third backtrack should exceed budget")
	}
}

func TestAssignmentSearchAssignAndRetreat(t *testing.T) {
	as := newAssignmentSearch(3)

	as.assign(2)
	if as.step != 1 {
		t.Fatalf("step after assign = %d, want 1", as.step)
	}
	if !as.usedAsReceiver[2] {
		t.Fatal("assigned receiver should be marked")
	}
	if as.assignment[0] != 2 {
		t.Fatalf("assignment[0] = %d, want 2", as.assignment[0])
	}

	as.assign(0)
	if as.step != 2 {
		t.Fatalf("step after second assign = %d, want 2", as.step)
	}

	ok := as.retreat()
	if !ok {
		t.Fatal("retreat should succeed")
	}
	// After retreat: current (step=2) is reset, previous (step=1) is undone.
	// We land back at step=1 to retry with a different candidate.
	if as.step != 1 {
		t.Fatalf("step after retreat = %d, want 1", as.step)
	}
	if as.usedAsReceiver[0] {
		t.Fatal("retreated receiver at step 1 should be unmarked")
	}
}

func TestAssignmentSearchIsComplete(t *testing.T) {
	as := newAssignmentSearch(2)
	if as.isComplete(2) {
		t.Fatal("should not be complete initially")
	}
	as.assign(1)
	as.assign(0)
	if !as.isComplete(2) {
		t.Fatal("should be complete after assigning all")
	}
}
