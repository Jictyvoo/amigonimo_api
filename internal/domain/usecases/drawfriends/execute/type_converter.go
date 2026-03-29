package execute

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/execute/matcher"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

// toMatcherParticipants converts domain participants into matcher DTOs,
// pre-computing the allowed receivers list (all others minus self minus denylist).
func toMatcherParticipants(participants []entities.Participant) []matcher.Participant {
	idSet := make(map[entities.HexID]struct{}, len(participants))
	for _, p := range participants {
		idSet[p.ID] = struct{}{}
	}

	result := make([]matcher.Participant, len(participants))
	for i, p := range participants {
		denySet := make(map[entities.HexID]struct{}, len(p.DenyList))
		for _, d := range p.DenyList {
			denySet[d.ID] = struct{}{}
		}

		allowed := make([]matcher.ParticipantID, 0, len(participants)-1)
		for _, other := range participants {
			if other.ID == p.ID {
				continue
			}
			if _, denied := denySet[other.ID]; denied {
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

// toPairingResults converts matcher pairings back to domain DrawResultItems.
func toPairingResults(
	pairings []matcher.Pairing,
	participants []entities.Participant,
) []entities.DrawResultItem {
	byID := make(map[entities.HexID]entities.Participant, len(participants))
	for _, p := range participants {
		byID[p.ID] = p
	}

	results := make([]entities.DrawResultItem, len(pairings))
	for i, pair := range pairings {
		results[i] = entities.DrawResultItem{
			Giver:    byID[entities.HexID(pair.GiverID)],
			Receiver: byID[entities.HexID(pair.ReceiverID)],
		}
	}
	return results
}
