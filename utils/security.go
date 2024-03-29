package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(hashFromDatabase string, password string) bool {
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(hashFromDatabase), []byte(password)); err != nil {
		return false
	}

	return true
}
