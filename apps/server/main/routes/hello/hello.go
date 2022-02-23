package hello

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Setup(g *echo.Group) {

	g.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

}
