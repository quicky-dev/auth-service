package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/db"
	"log"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Register a new user
func Register(c echo.Context) error {
	credentials := new(registerRequest)

	if err := c.Bind(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, authError{"The body doesn't match RegisterCredentials"})
	}

	if err := credentials.ValidateFields(); err != nil {
		return c.JSON(http.StatusBadRequest, authError{err.Error()})
	}

	user := credentials.ToUser()
	objectID, err := user.Save()

	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, authError{err.Error()})
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
		return c.JSON(http.StatusBadRequest, authError{"The body is malformed."})
	}

	if err := credentials.ValidateFields(); err != nil {
		return c.JSON(http.StatusBadRequest, authError{err.Error()})
	}

	// Convert the request to a user and then attempt to login with that user.
	user := credentials.ToUser()
	err := user.Login()

	if err != nil {
		return c.JSON(400, authError{"Login failed."})
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &jwtClaims{
		Username: user.Username,
		ID:       user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			authError{"Couldn't issue an access token."})
	}

	c.SetCookie(&http.Cookie{
		Name:    "access-token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.SetCookie(&http.Cookie{
		Name:    "refresh-token",
		Value:   user.RefreshToken,
		Expires: user.LastSignIn.Add(5 * 24 * time.Hour),
	})

	return c.JSON(http.StatusOK, loginResponse{user.Username})
}

// Refresh a users jwt token. This will keep users tokens in rotation and keep
// our users secure.
func RefreshToken(c echo.Context) error {
	accessCookie, err := c.Cookie("access-token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, authError{"The user isn't currently signed in"})
	}

	refreshCookie, err := c.Cookie("refresh-token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, authError{"The user doesn't have a refresh token."})
	}

	jwtString := accessCookie.Value
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(
		jwtString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err == jwt.ErrSignatureInvalid {
		return c.JSON(http.StatusUnauthorized, authError{"Invalid JWT Token."})
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, authError{"Bad request."})
	}

	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, authError{"Invalid JWT Token."})
	}

	user := claims.ToUser()
	if err := user.ValidateAndReplaceToken(refreshCookie.Value); err != nil {
		return c.JSON(http.StatusInternalServerError, authError{err.Error()})
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims = &jwtClaims{
		Username: user.Username,
		ID:       user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			authError{"Couldn't issue an access token."})
	}

	c.SetCookie(&http.Cookie{
		Name:    "access-token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.SetCookie(&http.Cookie{
		Name:    "refresh-token",
		Value:   user.RefreshToken,
		Expires: user.LastSignIn.Add(5 * 24 * time.Hour),
	})

	return c.JSON(http.StatusOK, refreshTokenResponse{user.Username})
}

// Verify a users token.
func VerifyToken(c echo.Context) error {
	accessCookie, err := c.Cookie("access-token")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, authError{"The user isn't currently signed in"})
	}
	return c.String(http.StatusOK, "Hello world")
}
