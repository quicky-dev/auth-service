package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/db"
	"github.com/quicky-dev/auth-service/util"
	"log"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
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

func formatVerificationURL(verificationCode string, userID string) string {
	if os.Getenv("PRODUCTION") == "" {
		return fmt.Sprintf("http://localhost:3000/verify?v=%s&u=%s", verificationCode, userID)
	} else {
		return fmt.Sprintf("https://auth.quicky.dev/verify?v=%s&u=%s", verificationCode, userID)
	}
}

// Register a new user
func Register(c echo.Context) error {
	credentials := new(RegisterRequest)

	if err := c.Bind(credentials); err != nil {
		return c.JSON(400, AuthError{"The body doesn't match RegisterCredentials"})
	}

	user, err := registerRequestToUnverifiedUser(credentials)

	if err != nil {
		log.Println(err.Error())
		return c.JSON(500, AuthError{err.Error()})
	}

	objectID, err := db.AddUnverifiedUser(user)

	if err != nil {
		log.Println(err.Error())
		return c.JSON(500, AuthError{err.Error()})
	}

	return c.String(http.StatusOK, formatVerificationURL(user.VerificationCode, objectID))
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
