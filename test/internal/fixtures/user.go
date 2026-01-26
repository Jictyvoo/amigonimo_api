package fixtures

import (
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
)

type UserBuilder struct {
	instance *genmodels.User
}

func NewUser() *UserBuilder {
	uid, err := uuid.NewV7()
	if err != nil {
		log.Panicf("failed to generate uuid: %s", err)
	}

	newBuilder := &UserBuilder{
		instance: &genmodels.User{
			ID:        uid[:],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Fullname:  "Test User " + uid.String(),
			Email:     "test-" + uid.String() + "@example.com",
			Username:  "user-" + uid.String(),
		},
	}

	return newBuilder.WithPassword("password")
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.instance.Email = email
	return b
}

func (b *UserBuilder) WithFullname(fullname string) *UserBuilder {
	b.instance.Fullname = fullname
	return b
}

func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.instance.Username = username
	return b
}

func (b *UserBuilder) WithPassword(rawPassword string) *UserBuilder {
	userEntity := entities.UserBasic{Password: rawPassword}
	encryptedPass, _ := userEntity.EncryptPassword()
	b.instance.Password = string(encryptedPass)
	return b
}

func (b *UserBuilder) Build() *genmodels.User {
	return b.instance
}
