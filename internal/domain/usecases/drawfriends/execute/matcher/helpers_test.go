package matcher

import "testing"

func TestSortByConstraint(t *testing.T) {
	a := Participant{ID: pid(1), AllowedReceivers: []ParticipantID{pid(2)}}
	b := Participant{ID: pid(2), AllowedReceivers: []ParticipantID{pid(1), pid(3)}}
	c := Participant{ID: pid(3), AllowedReceivers: []ParticipantID{pid(1), pid(2), pid(4)}}

	tests := []struct {
		name      string
		ascending bool
		wantFirst ParticipantID
		wantLast  ParticipantID
	}{
		{
			name:      "ascending sorts most constrained first",
			ascending: true,
			wantFirst: pid(1),
			wantLast:  pid(3),
		},
		{
			name:      "descending sorts least constrained first",
			ascending: false,
			wantFirst: pid(3),
			wantLast:  pid(1),
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				sorted := sortByConstraint([]Participant{c, a, b}, tt.ascending)
				if sorted[0].ID != tt.wantFirst {
					t.Fatalf("first = %v, want %v", sorted[0].ID, tt.wantFirst)
				}
				if sorted[len(sorted)-1].ID != tt.wantLast {
					t.Fatalf("last = %v, want %v", sorted[len(sorted)-1].ID, tt.wantLast)
				}
			},
		)
	}
}

func TestSortByConstraintTieBreaksByID(t *testing.T) {
	a := Participant{ID: pid(2), AllowedReceivers: []ParticipantID{pid(1)}}
	b := Participant{ID: pid(1), AllowedReceivers: []ParticipantID{pid(2)}}

	sorted := sortByConstraint([]Participant{a, b}, true)
	if sorted[0].ID != pid(1) {
		t.Fatalf("tie-break: first = %v, want %v", sorted[0].ID, pid(1))
	}
}

func TestBuildAllowedSet(t *testing.T) {
	p := Participant{
		ID:               pid(1),
		AllowedReceivers: []ParticipantID{pid(2), pid(3)},
	}
	set := buildAllowedSet(p)

	if _, ok := set[pid(2)]; !ok {
		t.Fatal("expected pid(2) in allowed set")
	}
	if _, ok := set[pid(3)]; !ok {
		t.Fatal("expected pid(3) in allowed set")
	}
	if _, ok := set[pid(1)]; ok {
		t.Fatal("pid(1) should not be in allowed set")
	}
}
