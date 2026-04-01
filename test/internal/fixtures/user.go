package fixtures

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/jictyvoo/amigonimo_api/internal/entities/authvalues"
	"github.com/jictyvoo/amigonimo_api/internal/infra/repositories/mysqlrepo/genmodels"
)

type UserBuilder struct {
	instance *genmodels.User
	profile  *genmodels.UserProfile
}

func NewUser() *UserBuilder {
	uid, err := uuid.NewV7()
	if err != nil {
		log.Panicf("failed to generate uuid: %s", err)
	}
	profileID, err := uuid.NewV7()
	if err != nil {
		log.Panicf("failed to generate profile uuid: %s", err)
	}
	now := time.Now()
	defaultFullname := "Test User " + uid.String()

	newBuilder := &UserBuilder{
		instance: &genmodels.User{
			ID:        uid[:],
			CreatedAt: now,
			UpdatedAt: now,
			Email:     "test-" + uid.String() + "@example.com",
			Username:  "user-" + uid.String(),
		},
		profile: &genmodels.UserProfile{
			ID:        profileID[:],
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    uid[:],
			Fullname:  sql.NullString{String: defaultFullname, Valid: true},
		},
	}

	return newBuilder.WithPassword("password")
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.instance.Email = email
	return b
}

func (b *UserBuilder) WithFullname(fullname string) *UserBuilder {
	b.profile.Fullname = sql.NullString{String: fullname, Valid: fullname != ""}
	return b
}

func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.instance.Username = username
	return b
}

func (b *UserBuilder) WithPassword(rawPassword string) *UserBuilder {
	userEntity := authvalues.UserBasic{Password: rawPassword}
	encryptedPass, _ := userEntity.EncryptPassword()
	b.instance.Password = string(encryptedPass)
	return b
}

func (b *UserBuilder) Build() *genmodels.User {
	return b.instance
}

func (b *UserBuilder) BuildProfile() *genmodels.UserProfile {
	return b.profile
}
