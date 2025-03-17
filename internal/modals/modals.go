package modals

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {

  passwordBytes := []byte(password) 
  hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

  if err != nil {
    return "", err
  }

  return string(hashedPassword), nil 

}

func generateSessionToken() string {
  b := make([]byte, 32)
  rand.Read(b)
  return base64.URLEncoding.EncodeToString(b)
}

func getTime(input string) (time.Time, error) {
  layout := "2006-01-02 15:04:05"  
  t, err := time.Parse(layout, input)
  if err != nil {
    fmt.Println("Error parsing time:", err)
    return time.Now(), err
  }
  return t, nil
}

