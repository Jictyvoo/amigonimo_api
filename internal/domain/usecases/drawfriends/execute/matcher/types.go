package matcher

// ParticipantID is a draw-local identifier decoupled from entities.HexID.
// It mirrors the 16-byte UUID representation for zero-cost conversion.
type ParticipantID [16]byte

func (participantID ParticipantID) Compare(other ParticipantID) int {
	for i := range participantID {
		if participantID[i] < other[i] {
			return -1
		}
		if participantID[i] > other[i] {
			return 1
		}
	}
	return 0
}

// Participant is the draw-local DTO carrying only what the algorithms need:
// the participant's identity and its pre-computed list of valid receivers.
type Participant struct {
	ID               ParticipantID
	AllowedReceivers []ParticipantID
}

// Pairing represents a single giver→receiver assignment produced by a draw strategy.
type Pairing struct {
	GiverID    ParticipantID
	ReceiverID ParticipantID
}
