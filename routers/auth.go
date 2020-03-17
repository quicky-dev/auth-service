package routers

import (
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/controllers"
)

func AttachAuth(app *echo.Echo) {
	app.POST("/auth/register", controllers.Register)
}
