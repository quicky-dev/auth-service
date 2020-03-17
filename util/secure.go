package util

import (
	"crypto/rand"
)

// These two functions have been taken from:
// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
// Essentially, we need secure random strings for verifying a users email.

func generateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func GenerateRandomString(size int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	randomBytes, err := generateRandomBytes(size)

	if err != nil {
		return "", err
	}

	// Iterate over the random bytes and convert into string bytes.
	for i, b := range randomBytes {
		randomBytes[i] = letters[b%byte(len(letters))]
	}

	return string(randomBytes), nil
}
