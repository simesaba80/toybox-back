package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type AssetController struct {
	assetUsecase usecase.IAssetUseCase
}

func NewAssetController(assetUsecase usecase.IAssetUseCase) *AssetController {
	return &AssetController{
		assetUsecase: assetUsecase,
	}
}

func (ac *AssetController) UploadAsset(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*schema.JWTCustomClaims)
	userID := claims.UserID
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "File is required")
	}
	assetURL, err := ac.assetUsecase.UploadFile(c.Request().Context(), file, userID)
	if err != nil {
		c.Logger().Error("Failed to upload asset: %w", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to upload asset")
	}
	return c.JSON(http.StatusOK, map[string]string{"url": *assetURL})
}
