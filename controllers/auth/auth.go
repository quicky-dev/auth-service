package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/db"
	"log"
	"net/http"
	"os"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	Username string `json:"username"`
	Uid      string `json:"uid"`
	jwt.StandardClaims
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

	objectID, err := db.CreateUnverifiedUser(user)

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
