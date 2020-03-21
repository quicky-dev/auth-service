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

type claims struct {
	Username string `json:"username"`
	Uid      string `json:"uid"`
	jwt.StandardClaims
}

// Register a new user
func Register(c echo.Context) error {
	credentials := new(registerRequest)

	if err := c.Bind(credentials); err != nil {
		return c.JSON(400, authError{"The body doesn't match RegisterCredentials"})
	}

	if err := credentials.ValidateFields(); err != nil {
		return c.JSON(400, authError{err.Error()})
	}

	user, err := credentials.ToUser()

	if err != nil {
		log.Println(err.Error())
		return c.JSON(500, authError{err.Error()})
	}

	objectID, err := user.Save()

	if err != nil {
		log.Println(err.Error())
		return c.JSON(500, authError{err.Error()})
	}

	return c.String(http.StatusOK, formatVerificationURL(user.VerificationCode, objectID))
}

func VerifyEmail(c echo.Context) error {
	verificationCode := c.QueryParam("v")
	userID := c.QueryParam("u")

	user := new(db.User)
	_, err := user.VerifyEmail(userID, verificationCode)

	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusForbidden, authError{err.Error()})
	}

	log.Printf("Verified account for: %s", user.Username)

	return c.JSON(http.StatusOK, verifyEmailResponse{user.Username})
}

// Login a pre existing user
func Login(c echo.Context) error {
	credentials := new(loginRequest)

	if err := c.Bind(credentials); err != nil {
		return c.JSON(400, authError{"The body is malformed."})
	}

	if err := credentials.ValidateFields(); err != nil {
		return c.JSON(400, authError{err.Error()})
	}

	user, err := credentials.ToUser()

	if err != nil {
		return c.JSON(400, authError{err.Error()})
	}

	err = user.Login()

	if err != nil {
		return c.JSON(400, authError{"Login failed."})
	}

	return c.JSON(http.StatusOK, loginResponse{user.Username})
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
