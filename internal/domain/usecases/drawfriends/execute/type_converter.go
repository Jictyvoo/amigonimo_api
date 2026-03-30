package execute

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/drawdto"
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/execute/matcher"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// toMatcherParticipants converts domain participants into matcher DTOs,
// pre-computing the allowed receivers list (all others minus self).
// TODO: enforce denylist constraints — currently all participants are allowed as receivers
func toMatcherParticipants(participants []entities.Participant) []matcher.Participant {
	result := make([]matcher.Participant, len(participants))
	for i, p := range participants {
		allowed := make([]matcher.ParticipantID, 0, len(participants)-1)
		for _, other := range participants {
			if other.ID == p.ID {
				continue
			}
			allowed = append(allowed, matcher.ParticipantID(other.ID))
		}
		result[i] = matcher.Participant{
			ID:               matcher.ParticipantID(p.ID),
			AllowedReceivers: allowed,
		}
	}
	return result
}

// toPairingResults converts matcher pairings back to drawdto.DrawResultItems.
func toPairingResults(
	pairings []matcher.Pairing,
	participants []entities.Participant,
) []drawdto.DrawResultItem {
	byID := make(map[entities.HexID]entities.Participant, len(participants))
	for _, p := range participants {
		byID[p.ID] = p
	}

	results := make([]drawdto.DrawResultItem, len(pairings))
	for i, pair := range pairings {
		results[i] = drawdto.DrawResultItem{
			GiverParticipantID:    entities.HexID(pair.GiverID),
			ReceiverParticipantID: entities.HexID(pair.ReceiverID),
		}
	}
	return results
}
