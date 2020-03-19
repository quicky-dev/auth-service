package auth

import (
	"errors"
	"github.com/quicky-dev/auth-service/db"
	"github.com/quicky-dev/auth-service/util"
)

// ---------------------------------- REQUESTS ---------------------------------

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (this *RegisterRequest) validateFields() error {
	if this.Username == "" {
		return errors.New("The username field has been left blank.")
	}

	if len(this.Username) <= 3 {
		return errors.New("The username is less than or equal to 3 characters long.")
	}

	if this.Password == "" {
		return errors.New("The password field has been left blank.")
	}

	if len(this.Password) <= 5 {
		return errors.New("The password is less than or equal to 5 characters.")
	}

	if this.Email == "" {
		return errors.New("The email field has been left blank.")
	}

	return nil
}

func (this *RegisterRequest) toUnverifiedUser() (*db.UnverifiedUser, error) {
	user := new(db.UnverifiedUser)
	user.Username = this.Username
	user.Email = this.Email
	hashedPassword, err := util.HashAndSaltPassword(this.Password)

	if err != nil {
		return &db.UnverifiedUser{}, err
	}

	user.Password = hashedPassword

	generatedString, err := util.GenerateRandomString(64)
	if err != nil {
		return &db.UnverifiedUser{}, err
	}

	user.VerificationCode = generatedString

	return user, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
