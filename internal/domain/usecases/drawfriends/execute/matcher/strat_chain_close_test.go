package matcher

import "testing"

func TestChainCloseStrategy(t *testing.T) {
	tests := baseTestCases()
	runStrategyTests(t, ChainCloseStrategy{}, tests)
}

func TestChainCloseProducesSingleCycle(t *testing.T) {
	tests := []struct {
		name                 string
		numberOfParticipants int
		denyPairs            [][2]byte
	}{
		{name: "3 participants", numberOfParticipants: 3},
		{name: "5 participants", numberOfParticipants: 5},
		{name: "6 with denylist", numberOfParticipants: 6, denyPairs: [][2]byte{{1, 2}}},
		{
			name:                 "8 with cross denylists",
			numberOfParticipants: 8,
			denyPairs:            [][2]byte{{1, 2}, {2, 1}, {3, 4}},
		},
	}

	s := ChainCloseStrategy{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			participants := makeParticipants(tt.numberOfParticipants, tt.denyPairs)
			pairs, err := s.Execute(participants)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			validatePairings(t, pairs, participants)
			if !isSingleCycle(pairs) {
				t.Fatal("expected single cycle, got multiple cycles")
			}
		})
	}
}
