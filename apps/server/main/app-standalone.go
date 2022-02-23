//go:build !appengine && !appenginevm
// +build !appengine,!appenginevm

package main

import (
	"os"
	"server/generated/db"
	"server/main/routes/hello"
	"server/main/routes/users"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func createMux() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	return e
}

func main() {
	db.Connect()

	e := createMux()
	g := e.Group("/api")
	users.Setup(g)
	hello.Setup(g)

	if os.Getenv("GIN_MODE") == "release" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		e.Static("/", "webapp")
		e.Start(":" + port)

	} else {
		e.Start(":4000")

	}

}
