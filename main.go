package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/routers"
)

// Main entry point for the auth service.
func main() {
	app := echo.New()

	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	routers.AttachAuth(app)

	app.Logger.Fatal(app.Start(":8081"))
}
