package router

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/simesaba80/toybox-back/internal/infrastructure/config"
	customejwt "github.com/simesaba80/toybox-back/internal/infrastructure/external/custome-jwt"
	"github.com/simesaba80/toybox-back/internal/interface/controller"
	"github.com/simesaba80/toybox-back/pkg/echovalidator"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	echo              *echo.Echo
	UserController    *controller.UserController
	WorkController    *controller.WorkController
	CommentController *controller.CommentController
	AuthController    *controller.AuthController
}

func NewRouter(e *echo.Echo, uc *controller.UserController, wc *controller.WorkController, cc *controller.CommentController, ac *controller.AuthController) *Router {
	return &Router{
		echo:              e,
		UserController:    uc,
		WorkController:    wc,
		CommentController: cc,
		AuthController:    ac,
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

	r.echo.GET("/auth/discord", r.AuthController.GetDiscordAuthURL)
	r.echo.GET("/auth/discord/callback", r.AuthController.AuthenticateUser)
	r.echo.POST("/auth/refresh", r.AuthController.RegenerateToken)
	r.echo.POST("/users", r.UserController.CreateUser)
	r.echo.GET("/users", r.UserController.GetAllUsers)
	r.echo.GET("/users/auth", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	r.echo.POST("/works", r.WorkController.CreateWork)
	r.echo.GET("/works", r.WorkController.GetAllWorks)
	r.echo.GET("/works/:work_id", r.WorkController.GetWorkByID)

	r.echo.GET("/works/:work_id/comments", r.CommentController.GetCommentsByWorkID)
	r.echo.POST("/works/:work_id/comments", r.CommentController.CreateComment)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(customejwt.JWTCustomClaims)
		},
		SigningKey: []byte(config.TOKEN_SECRET),
	}
	e := r.echo.Group("/works", echojwt.WithConfig(config))
	e.Use(echojwt.WithConfig(config))
	e.GET("/", r.WorkController.GetAllWorks)
	return r.echo
}
