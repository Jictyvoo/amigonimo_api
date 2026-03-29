package matcher

import "testing"

func TestGreedyStrategy(t *testing.T) {
	tests := baseTestCases()
	runStrategyTests(t, GreedyStrategy{}, tests)
}

func TestGreedyWithHeavyDenylists(t *testing.T) {
	// Each participant denies one other — still solvable.
	participants := makeParticipants(6, [][2]byte{
		{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 1},
	})
	s := GreedyStrategy{}
	pairs, err := s.Execute(participants)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	validatePairings(t, pairs, participants)
}
