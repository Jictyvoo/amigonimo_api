package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	testCases := []struct {
		name         string
		participant1 ParticipantID
		participant2 ParticipantID
		expected     int
	}{
		{
			name:         "Equal participants",
			participant1: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			participant2: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected:     0,
		},
		{
			name:         "Participant1 greater than Participant2",
			participant1: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 17},
			participant2: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected:     1,
		},
		{
			name:         "Different lengths",
			participant1: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			participant2: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			expected:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				result := tc.participant1.Compare(tc.participant2)
				assert.Equal(t, tc.expected, result, "Compare returned incorrect result")
			},
		)
	}
}
