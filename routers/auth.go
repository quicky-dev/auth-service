package routers

import (
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/controllers/auth"
)

// AttachAuth attaches all auth endpoints to the main echo instance
func RegisterAuthRoutes(app *echo.Echo) {
	// GET
	app.GET("/verify/email", auth.VerifyEmail)
	app.GET("/refresh", auth.RefreshToken)

	// POST
	app.POST("/register", auth.Register)
	app.POST("/login", auth.Login)
	app.POST("/verify/token", auth.VerifyToken)
}
