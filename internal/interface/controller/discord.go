package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type DiscordController struct {
	discordUsecase *usecase.DiscordUsecase
}

func NewDiscordController(discordUsecase *usecase.DiscordUsecase) *DiscordController {
	return &DiscordController{
		discordUsecase: discordUsecase,
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
func (dc *DiscordController) GetDiscordAuthURL(c echo.Context) error {
	authURL, err := dc.discordUsecase.GetDiscordAuthURL(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, schema.ToGetDiscordAuthURLResponse(authURL))
}
