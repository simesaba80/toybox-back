package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/pkg/echovalidator"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	echo           *echo.Echo
	UserController *controller.UserController
	WorkController *controller.WorkController
}

func NewRouter(e *echo.Echo, uc *controller.UserController, wc *controller.WorkController) *Router {
	return &Router{
		echo:           e,
		UserController: uc,
		WorkController: wc,
	}
}

func (r *Router) Setup() *echo.Echo {
	r.echo.Validator = echovalidator.NewValidator()
	r.echo.Use(middleware.Logger())
	r.echo.Use(middleware.Recover())
	r.echo.Use(middleware.CORS())

	r.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	r.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	r.echo.POST("/users", r.UserController.CreateUser)
	r.echo.GET("/users", r.UserController.GetAllUsers)

	r.echo.POST("/works", r.WorkController.CreateWork)
	r.echo.GET("/works", r.WorkController.GetAllWorks)
	r.echo.GET("/works/:work_id", r.WorkController.GetWorkByID)

	return r.echo
}
