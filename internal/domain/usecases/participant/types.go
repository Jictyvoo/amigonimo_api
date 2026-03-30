package participant

import "github.com/jictyvoo/amigonimo_api/internal/entities"

// Summary is the read-model DTO returned when listing participants.
// It embeds the full Participant entity and adds FullName,
// which requires a JOIN on user_profiles and is not part of the core entity.
type Summary struct {
	entities.Participant
	FullName string
}
