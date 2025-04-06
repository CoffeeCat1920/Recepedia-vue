package modals

import (
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Password string
}


var (
  admin *Admin
  admin_password = "123" 
) 


func NewAdmin() (*Admin) {
	if admin != nil {
		return admin
	}

	hashedPassword,_ := hashPassword(admin_password) 

	admin = &Admin{
		Password: hashedPassword,
	}

	return admin
}

func (admin *Admin) CheckPassword(password string) bool {
	hashedPassword := []byte(admin.Password)
	passwordBytes := []byte(password)
	
	err := bcrypt.CompareHashAndPassword(hashedPassword, passwordBytes)

	return err == nil
}
