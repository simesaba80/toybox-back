package controller

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
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
		return handleAssetError(c, err)
	}
	return c.JSON(http.StatusOK, schema.ToUploadAssetResponse(*assetURL))
}

func handleAssetError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domainerrors.ErrInvalidRequestBody):
		return echo.NewHTTPError(http.StatusBadRequest, "無効なリクエストです")
	case errors.Is(err, domainerrors.ErrFailedToOpenFile):
		return echo.NewHTTPError(http.StatusInternalServerError, "ファイルの読み込みに失敗しました")
	case errors.Is(err, domainerrors.ErrFailedToUploadFile):
		return echo.NewHTTPError(http.StatusInternalServerError, "ファイルのアップロードに失敗しました")
	case errors.Is(err, domainerrors.ErrFailedToCreateAsset):
		return echo.NewHTTPError(http.StatusInternalServerError, "アセットの作成に失敗しました")
	}
	c.Logger().Error("Failed to upload asset: %w", err)
	return echo.NewHTTPError(http.StatusInternalServerError, "サーバーエラーが発生しました")
}
