package entities

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	UserBasic struct {
		Username string
		Email    string
		Password string
	}
	User struct {
		ID            HexID
		FullName      string
		VerifiedAt    time.Time
		RememberToken string
		UserBasic
	}
)

func (ub UserBasic) ObfuscateEmail() string {
	obfuscate := strings.Builder{}
	firstHalf, hostHalf, _ := strings.Cut(ub.Email, "@")
	for index, character := range []byte(firstHalf) {
		if index <= (len(firstHalf) / 3) {
			obfuscate.WriteByte(character)
		} else if index%2 == 0 {
			obfuscate.WriteRune('*')
		}
	}
	obfuscate.WriteRune('@')
	obfuscate.WriteString(hostHalf)
	return obfuscate.String()
}

// EncryptPassword hashes a password using bcrypt with the configured cost.
func (ub UserBasic) EncryptPassword() ([]byte, error) {
	const passwordCost = 11
	return bcrypt.GenerateFromPassword([]byte(ub.Password), passwordCost)
}

// ComparePassword checks the password hash against another password.
// Is expected that the UserBasic stored password is already encrypted by the EncryptPassword method,
// and the other password should not be encrypted.
func (ub UserBasic) ComparePassword(otherPass string) (ok bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(ub.Password), []byte(otherPass))
	ok = err == nil
	return ok, err
}
