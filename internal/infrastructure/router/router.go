package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
)

type Router struct {
	echo           *echo.Echo
	UserController *controller.UserController
}

func NewRouter(e *echo.Echo, uc *controller.UserController) *Router {
	return &Router{
		echo:           e,
		UserController: uc,
	}
}

func (r *Router) Setup() *echo.Echo {
	r.echo.Use(middleware.Logger())
	r.echo.Use(middleware.Recover())
	r.echo.Use(middleware.CORS())

	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	r.echo.POST("/users", r.UserController.CreateUser)
	r.echo.GET("/users", r.UserController.GetAllUsers)

	return r.echo
}
