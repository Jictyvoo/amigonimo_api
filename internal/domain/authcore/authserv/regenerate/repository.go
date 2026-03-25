package regenerate

import "github.com/jictyvoo/amigonimo_api/internal/entities"

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=repository_mock_test.go -package=regenerate github.com/jictyvoo/amigonimo_api/internal/domain/authcore/authserv/regenerate Repository

type Repository interface {
	CheckAuthenticationByRefreshToken(authToken string) (entities.AuthenticationToken, error)
	UpsertAuthToken(authentication *entities.AuthenticationToken) error
}
