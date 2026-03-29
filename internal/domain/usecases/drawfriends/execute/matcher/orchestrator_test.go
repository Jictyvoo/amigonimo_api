package matcher

import "testing"

func TestOrchestratorExecute(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		denyPairs [][2]byte
		wantErr   bool
	}{
		{
			name: "basic 3 participants",
			n:    3,
		},
		{
			name: "5 participants no denylists",
			n:    5,
		},
		{
			name: "6 participants with denylists",
			n:    6,
			denyPairs: [][2]byte{
				{1, 2},
				{2, 1},
				{3, 4},
			},
		},
		{
			name:    "too few participants",
			n:       2,
			wantErr: true,
		},
		{
			name: "4 participants with cross denylists",
			n:    4,
			denyPairs: [][2]byte{
				{1, 2},
				{2, 1},
				{3, 4},
				{4, 3},
			},
		},
	}

	orch := NewOrchestrator()

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				participants := makeParticipants(tt.n, tt.denyPairs)
				pairs, err := orch.Execute(participants)

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
			},
		)
	}
}
