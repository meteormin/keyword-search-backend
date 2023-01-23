package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(fromPassword), err
}

func HashCheck(hashPass string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}
