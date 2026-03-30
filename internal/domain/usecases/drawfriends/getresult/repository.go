package getresult

import (
	"github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/drawdto"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=getresult github.com/jictyvoo/amigonimo_api/internal/domain/usecases/drawfriends/getresult Repository

type Repository interface {
	GetDrawResultForUser(secretFriendID, userID entities.HexID) (drawdto.DrawResultItem, error)
}
