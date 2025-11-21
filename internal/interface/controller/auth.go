package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type AuthController struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthController(authUsecase *usecase.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

// GetDiscordAuthURL godoc
// @Summary Get Discord authentication URL
// @Description Get Discord authentication URL
// @Tags auth
// @Produce json
// @Success 200 {object} schema.GetDiscordAuthURLResponse
// @Failure 500 {object} echo.HTTPError
// @Router /auth/discord/ [get]
func (ac *AuthController) GetDiscordAuthURL(c echo.Context) error {
	authURL, err := ac.authUsecase.GetDiscordAuthURL(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, schema.ToGetDiscordAuthURLResponse(authURL))
}

// GetDiscordToken godoc
// @Summary Get Discord token by code
// @Description Get Discord token by code
// @Tags auth
// @Produce json
// @Success 200 {object} schema.GetDiscordTokenResponse
// @Failure 500 {object} echo.HTTPError
// @Router /auth/discord/callback [get]
// @Param code query string true "Discord code"
func (ac *AuthController) AuthenticateUser(c echo.Context) error {
	code := c.QueryParam("code")
	appToken, refreshToken, err := ac.authUsecase.AuthenticateUser(c.Request().Context(), code)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, schema.ToGetDiscordTokenResponse(appToken, refreshToken))
}

// RegenerateToken godoc
// @Summary Regenerate token
// @Description Regenerate token
// @Tags auth
// @Produce json
// @Success 200 {object} schema.GetDiscordTokenResponse
// @Failure 500 {object} echo.HTTPError
// @Router /auth/refresh [post]
// @Param refresh_token query string true "Refresh token"
func (ac *AuthController) RegenerateToken(c echo.Context) error {
	refreshToken := c.QueryParam("refresh_token")
	appToken, err := ac.authUsecase.RegenerateToken(c.Request().Context(), refreshToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, schema.ToRegenerateTokenResponse(appToken))
}
