package routers

import (
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/controllers"
)

// AttachAuth attaches all auth endpoints to the main echo instance
func RegisterAuthRoutes(app *echo.Echo) {
	app.POST("/auth/register", controllers.Register)
}
