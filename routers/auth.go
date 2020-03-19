package routers

import (
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/controllers/auth"
)

// AttachAuth attaches all auth endpoints to the main echo instance
func RegisterAuthRoutes(app *echo.Echo) {
	app.POST("/register", auth.Register)
	app.GET("/verify/email", auth.VerifyEmail)
}
