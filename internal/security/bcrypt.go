package security

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashed), err
}

func ComparePassword(hasedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))
	return err
}
