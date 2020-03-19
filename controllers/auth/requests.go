package auth

import (
	"github.com/quicky-dev/auth-service/db"
	"github.com/quicky-dev/auth-service/util"
)

// ---------------------------------- REQUESTS ---------------------------------

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// -------------------------------- CONVERSIONS --------------------------------

func registerRequestToUnverifiedUser(request *RegisterRequest) (*db.UnverifiedUser, error) {
	user := new(db.UnverifiedUser)
	user.Username = request.Username
	user.Email = request.Email
	hashedPassword, err := util.HashAndSaltPassword(request.Password)

	if err != nil {
		return new(db.UnverifiedUser), err
	}

	user.Password = hashedPassword

	generatedString, err := util.GenerateRandomString(64)

	if err != nil {
		return new(db.UnverifiedUser), err
	}

	user.VerificationCode = generatedString

	return user, nil
}
