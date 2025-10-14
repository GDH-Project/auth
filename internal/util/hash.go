package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	if hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return "", err
	} else {
		return string(hashed), nil
	}
}

func CheckPasswordHash(password, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false, err
	} else {
		return true, nil
	}
}
