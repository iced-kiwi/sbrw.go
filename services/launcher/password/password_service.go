package password_service

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUppercase := false
	hasLowercase := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsLower(char) {
			hasLowercase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if hasUppercase && hasLowercase && hasDigit {
			break
		}
	}
	return hasUppercase && hasLowercase && hasDigit
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
