package main

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/simesaba80/toybox-back/pkg/config"
	"github.com/simesaba80/toybox-back/pkg/db"
)

func main() {
	config.LoadEnv()
	db.Init()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
