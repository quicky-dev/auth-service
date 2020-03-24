package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/quicky-dev/auth-service/db"
)

// ---------------------------------- REGISTER ---------------------------------

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (this *registerRequest) ValidateFields() error {
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

func (this *registerRequest) ToUser() *db.User {
	user := new(db.User)
	user.Username = this.Username
	user.Email = this.Email
	user.Password = this.Password

	return user
}

// ------------------------------------ LOGIN ----------------------------------

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this *loginRequest) ValidateFields() error {
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

	return nil
}

func (this *loginRequest) ToUser() *db.User {
	user := new(db.User)

	user.Username = this.Username
	user.Password = this.Password

	return user
}

// ----------------------------------- CLAIMS ----------------------------------

type jwtClaims struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	jwt.StandardClaims
}

func (this *jwtClaims) ToUser() *db.User {
	user := new(db.User)
	user.Username = this.Username
	return user
}
