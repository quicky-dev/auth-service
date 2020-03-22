package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type authError struct {
	ErrorMsg string `json:"errorMessage"`
}

type verifyEmailResponse struct {
	Username string `json:"username"`
}

type loginResponse struct {
	Username string `json:"username"`
}

type claims struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	jwt.StandardClaims
}
