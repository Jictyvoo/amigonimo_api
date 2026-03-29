package matcher

import "testing"

func TestReverseGreedyStrategy(t *testing.T) {
	tests := baseTestCases()
	runStrategyTests(t, ReverseGreedyStrategy{}, tests)
}

func TestReverseGreedyWithHeavyDenylists(t *testing.T) {
	participants := makeParticipants(6, [][2]byte{
		{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 1},
	})
	s := ReverseGreedyStrategy{}
	pairs, err := s.Execute(participants)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	validatePairings(t, pairs, participants)
}
