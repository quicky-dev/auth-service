package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/util"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterCredentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Uid      string `json:"uid"`
	jwt.StandardClaims
}

type AuthError struct {
	ErrorMsg string `json:"errorMessage"`
}

// Register a new user
func Register(c echo.Context) error {
	credentials := new(RegisterCredentials)

	if err := c.Bind(credentials); err != nil {
		return c.JSON(400, AuthError{"The body doesn't match RegisterCredentials"})
	}

	generatedString, err := util.GenerateRandomString(64)

	if err != nil {
		return c.JSON(500, AuthError{"Couldn't generate a secure string."})
	}

	return c.String(http.StatusOK, generatedString)
}

// Login a pre existing user
func Login(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, world!")
}

// Refresh a users jwt token. This will keep users tokens in rotation and keep
// our users secure.
func RefreshToken(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, world!")
}

// Verify a users token.
func VerifyToken(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}
