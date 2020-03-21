package util

import (
	"golang.org/x/crypto/bcrypt"
	"log"
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

// TODO: Differentiate the type of possible errors that could arise from
// comparing a password. I.E., being able to
// For example,
func ComparePasswordWithHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
