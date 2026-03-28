package matcher

import "testing"

func TestConstraintGraphUnusedCandidates(t *testing.T) {
	participants := makeParticipants(4, [][2]byte{{1, 2}})
	graph := newConstraintGraph(participants, true)

	used := make([]bool, graph.size())

	// Participant 1 (pid(1)) denies pid(2), so candidates should not include pid(2)'s index.
	idx1 := graph.idxByID[pid(1)]
	candidates := graph.unusedCandidates(idx1, used)

	for _, cIdx := range candidates {
		if graph.sorted[cIdx].ID == pid(2) {
			t.Fatal("candidate list should not include denied participant")
		}
	}

	// Mark one as used — should not appear.
	idx3 := graph.idxByID[pid(3)]
	used[idx3] = true
	candidates = graph.unusedCandidates(idx1, used)
	for _, cIdx := range candidates {
		if cIdx == idx3 {
			t.Fatal("candidate list should not include used participant")
		}
	}
}

func TestConstraintGraphCanReceive(t *testing.T) {
	participants := makeParticipants(3, [][2]byte{{1, 2}})
	graph := newConstraintGraph(participants, true)

	idx1 := graph.idxByID[pid(1)]

	if graph.canReceive(idx1, pid(2)) {
		t.Fatal("pid(1) should not be able to receive pid(2) (denied)")
	}
	if !graph.canReceive(idx1, pid(3)) {
		t.Fatal("pid(1) should be able to receive pid(3)")
	}
}

func TestConstraintGraphAssignmentToPairings(t *testing.T) {
	participants := makeParticipants(3, nil)
	graph := newConstraintGraph(participants, true)

	// assignment[0]=1, assignment[1]=2, assignment[2]=0
	assignment := []int{1, 2, 0}
	pairs := graph.assignmentToPairings(assignment)

	if len(pairs) != 3 {
		t.Fatalf("got %d pairs, want 3", len(pairs))
	}
	for i, p := range pairs {
		if p.GiverID != graph.sorted[i].ID {
			t.Fatalf("pair[%d].GiverID = %v, want %v", i, p.GiverID, graph.sorted[i].ID)
		}
		if p.ReceiverID != graph.sorted[assignment[i]].ID {
			t.Fatalf(
				"pair[%d].ReceiverID = %v, want %v",
				i,
				p.ReceiverID,
				graph.sorted[assignment[i]].ID,
			)
		}
	}
}
