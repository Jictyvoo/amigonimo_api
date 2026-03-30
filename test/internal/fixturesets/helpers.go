package fixturesets

import "github.com/jictyvoo/amigonimo_api/internal/entities"

func mustHexIDFromBytes(raw []byte) entities.HexID {
	id, err := entities.NewHexIDFromBytes(raw)
	if err != nil {
		panic(err)
	}

	return id
}
