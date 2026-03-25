package session

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/jictyvoo/amigonimo_api/internal/domain/authcore/autherrs"
	"github.com/jictyvoo/amigonimo_api/internal/entities"
)

//go:generate go tool -modfile=../../../../../build/tools/go.mod mockgen -destination=session_mock_test.go -package=session github.com/jictyvoo/amigonimo_api/internal/domain/authcore/userserv/session Repository

type Repository interface {
	GetUserByAuthCode(authToken string) (entities.User, error)
}

func FindAndCheckUser(
	repo Repository,
	authToken string,
	password string,
) (user entities.User, err error) {
	user, err = repo.GetUserByAuthCode(authToken)
	if err != nil || user.ID.IsEmpty() {
		return entities.User{}, autherrs.ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return entities.User{}, autherrs.ErrWrongPassword
	}

	return user, nil
}
