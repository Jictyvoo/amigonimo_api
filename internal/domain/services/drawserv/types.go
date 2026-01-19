package drawserv

import "github.com/jictyvoo/amigonimo_api/internal/entities"

type DrawInput struct {
	Participants []entities.Participant
}

type DrawOutput struct {
	Pairs []entities.DrawResultItem
}
