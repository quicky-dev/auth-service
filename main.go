package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/quicky-dev/auth-service/routers"
)

func main() {
	log.SetPrefix("[auth] ")

	app := echo.New()

	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	routers.RegisterAuthRoutes(app)

	app.Logger.Fatal(app.Start(os.Getenv("PORT")))
}
