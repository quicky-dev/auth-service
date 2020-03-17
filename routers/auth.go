package routers

import (
	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/controllers"
)

func AttachAuth(app *echo.Echo) {
	authRouter := echo.NewRouter(app)
	authRouter.Add(echo.POST, "/auth/register", controllers.Register)
}
