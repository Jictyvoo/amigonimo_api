package matcher

import "slices"

// sortByConstraint returns a copy of participants sorted by len(AllowedReceivers).
// Ascending = true sorts most-constrained first. Ties are broken by ID for determinism.
func sortByConstraint(participants []Participant, ascending bool) []Participant {
	sorted := make([]Participant, len(participants))
	copy(sorted, participants)
	slices.SortFunc(
		sorted, func(a, b Participant) int {
			diff := len(a.AllowedReceivers) - len(b.AllowedReceivers)
			if !ascending {
				diff = -diff
			}
			if diff != 0 {
				return diff
			}
			return ParticipantID.Compare(a.ID, b.ID)
		},
	)
	return sorted
}

// buildAllowedSet creates a quick-lookup set of allowed receivers for a participant.
func buildAllowedSet(p Participant) map[ParticipantID]struct{} {
	m := make(map[ParticipantID]struct{}, len(p.AllowedReceivers))
	for _, r := range p.AllowedReceivers {
		m[r] = struct{}{}
	}
	return m
}
