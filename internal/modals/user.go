package modals

import (
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewUser(name, password string) *User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal("\nCan't hash the password for user\n")
		return nil
	}

	return &User{
		UUID:     uuid.NewString(),
		Name:     name,
		Password: string(hashedPassword),
	}
}

func (user *User) CheckPassword(password string) bool {
	hashedPassword := []byte(user.Password)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashedPassword, passwordBytes)
	return err == nil
}
