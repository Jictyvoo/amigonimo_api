package matcher

import "testing"

// strategyTestCase defines a reusable test scenario for any DrawStrategy.
type strategyTestCase struct {
	name                 string
	numberOfParticipants int
	denyPairs            [][2]byte
	wantErr              bool
}

// baseTestCases returns the common set of test scenarios that every strategy
// must handle. Individual strategy tests may append or modify this list.
func baseTestCases() []strategyTestCase {
	return []strategyTestCase{
		{
			name:                 "basic 3 participants no denylists",
			numberOfParticipants: 3,
		},
		{
			name:                 "4 participants no denylists",
			numberOfParticipants: 4,
		},
		{
			name:                 "5 participants no denylists",
			numberOfParticipants: 5,
		},
		{
			name:                 "6 participants with light denylist",
			numberOfParticipants: 6,
			denyPairs:            [][2]byte{{1, 2}, {3, 4}},
		},
		{
			name:                 "4 participants with cross denylists",
			numberOfParticipants: 4,
			denyPairs:            [][2]byte{{1, 2}, {2, 1}},
		},
		{
			name:                 "too few participants",
			numberOfParticipants: 2,
			wantErr:              true,
		},
	}
}

// runStrategyTests executes all test cases against a single strategy,
// calling validatePairings on successful results.
func runStrategyTests(t *testing.T, strategy DrawStrategy, tests []strategyTestCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			participants := makeParticipants(tt.numberOfParticipants, tt.denyPairs)
			pairs, err := strategy.Execute(participants)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Execute() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			validatePairings(t, pairs, participants)
		})
	}
}

func pid(b byte) ParticipantID {
	var id ParticipantID
	id[0] = b
	return id
}

func makeParticipants(n int, denyPairs [][2]byte) []Participant {
	ids := make([]ParticipantID, n)
	for i := range n {
		ids[i] = pid(byte(i + 1))
	}

	denyMap := make(map[ParticipantID]map[ParticipantID]struct{})
	for _, dp := range denyPairs {
		giver := pid(dp[0])
		denied := pid(dp[1])
		if denyMap[giver] == nil {
			denyMap[giver] = make(map[ParticipantID]struct{})
		}
		denyMap[giver][denied] = struct{}{}
	}

	participants := make([]Participant, n)
	for i, id := range ids {
		var allowed []ParticipantID
		for _, other := range ids {
			if other == id {
				continue
			}
			if dm, ok := denyMap[id]; ok {
				if _, denied := dm[other]; denied {
					continue
				}
			}
			allowed = append(allowed, other)
		}
		participants[i] = Participant{ID: id, AllowedReceivers: allowed}
	}
	return participants
}

// validatePairings checks structural properties common to all valid draw results.
func validatePairings(t *testing.T, pairs []Pairing, participants []Participant) {
	t.Helper()
	n := len(participants)

	if len(pairs) != n {
		t.Fatalf("got %d pairs, want %d", len(pairs), n)
	}

	ids := make(map[ParticipantID]struct{}, n)
	for _, p := range participants {
		ids[p.ID] = struct{}{}
	}

	givers := make(map[ParticipantID]struct{}, n)
	receivers := make(map[ParticipantID]struct{}, n)
	allowedSets := make(map[ParticipantID]map[ParticipantID]struct{}, n)
	for _, p := range participants {
		allowedSets[p.ID] = buildAllowedSet(p)
	}

	for _, pair := range pairs {
		if pair.GiverID == pair.ReceiverID {
			t.Fatalf("self-draw: %v gives to self", pair.GiverID)
		}
		if _, ok := ids[pair.GiverID]; !ok {
			t.Fatalf("unknown giver: %v", pair.GiverID)
		}
		if _, ok := ids[pair.ReceiverID]; !ok {
			t.Fatalf("unknown receiver: %v", pair.ReceiverID)
		}
		if _, ok := allowedSets[pair.GiverID][pair.ReceiverID]; !ok {
			t.Fatalf("disallowed pair: %v → %v", pair.GiverID, pair.ReceiverID)
		}
		if _, dup := givers[pair.GiverID]; dup {
			t.Fatalf("duplicate giver: %v", pair.GiverID)
		}
		if _, dup := receivers[pair.ReceiverID]; dup {
			t.Fatalf("duplicate receiver: %v", pair.ReceiverID)
		}
		givers[pair.GiverID] = struct{}{}
		receivers[pair.ReceiverID] = struct{}{}
	}

	if len(givers) != n {
		t.Fatalf("only %d unique givers, want %d", len(givers), n)
	}
	if len(receivers) != n {
		t.Fatalf("only %d unique receivers, want %d", len(receivers), n)
	}
}

// isSingleCycle checks whether the pairings form exactly one cycle.
func isSingleCycle(pairs []Pairing) bool {
	if len(pairs) == 0 {
		return false
	}
	next := make(map[ParticipantID]ParticipantID, len(pairs))
	for _, p := range pairs {
		next[p.GiverID] = p.ReceiverID
	}
	start := pairs[0].GiverID
	current := next[start]
	visited := 1
	for current != start && visited <= len(pairs) {
		current = next[current]
		visited++
	}
	return visited == len(pairs) && current == start
}

func TestStrategyPriorities(t *testing.T) {
	tests := []struct {
		strategy DrawStrategy
		want     Priority
	}{
		{GreedyStrategy{}, PriorityGreedy},
	}

	for _, tt := range tests {
		if got := tt.strategy.ResultPriority(); got != tt.want {
			t.Errorf("%T.ResultPriority() = %d, want %d", tt.strategy, got, tt.want)
		}
	}
}
