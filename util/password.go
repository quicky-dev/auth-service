package util

import (
	"golang.org/x/crypto/bcrypt"
)

// TODO: Experiment with the cost value to see if we can increase security
// while maintaining speed.
func HashAndSaltPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func ComparePasswordWithHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return false
	}

	return true
}
