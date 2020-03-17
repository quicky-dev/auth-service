package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// Main entry point for the auth service.
func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8081"))
}
