package controller

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type FavoriteController struct {
	favoriteUsecase usecase.IFavoriteUsecase
}

func NewFavoriteController(favoriteUsecase usecase.IFavoriteUsecase) *FavoriteController {
	return &FavoriteController{favoriteUsecase: favoriteUsecase}
}

// CreateFavorite godoc
// @Summary Create a favorite
// @Description Create a favorite
// @Tags favorites
// @Accept json
// @Produce json
// @Param work_id path string true "Work ID"
// @Success 201
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security BearerAuth
// @Router /auth/works/:work_id/favorite [post]
func (fc *FavoriteController) CreateFavorite(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*schema.JWTCustomClaims)
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	workIDStr := c.Param("work_id")
	workID, err := uuid.Parse(workIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work ID format")
	}
	err = fc.favoriteUsecase.CreateFavorite(c.Request().Context(), workID, userID)
	if err != nil {
		return handleFavoriteError(c, err)
	}
	return c.NoContent(http.StatusCreated)
}

// DeleteFavorite godoc
// @Summary Delete a favorite
// @Description Delete a favorite
// @Tags favorites
// @Accept json
// @Produce json
// @Param work_id path string true "Work ID"
// @Success 204
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security BearerAuth
// @Router /auth/works/:work_id/favorite [delete]
func (fc *FavoriteController) DeleteFavorite(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*schema.JWTCustomClaims)
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	workIDStr := c.Param("work_id")
	workID, err := uuid.Parse(workIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work ID format")
	}
	err = fc.favoriteUsecase.DeleteFavorite(c.Request().Context(), workID, userID)
	if err != nil {
		return handleFavoriteError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// CountFavoritesByWorkID godoc
// @Summary Count favorites by work ID
// @Description Count favorites by work ID
// @Tags favorites
// @Accept json
// @Produce json
// @Param work_id path string true "Work ID"
// @Success 200 {object} schema.CountFavoritesByWorkIDResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /works/:work_id/favorite [get]
func (fc *FavoriteController) CountFavoritesByWorkID(c echo.Context) error {
	workIDStr := c.Param("work_id")
	workID, err := uuid.Parse(workIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work ID format")
	}
	total, err := fc.favoriteUsecase.CountFavoritesByWorkID(c.Request().Context(), workID)
	if err != nil {
		return handleFavoriteError(c, err)
	}
	return c.JSON(http.StatusOK, schema.CountFavoritesByWorkIDResponse{Total: total})
}

// IsFavorite godoc
// @Summary Check if a user has favorited a work
// @Description Check if a user has favorited a work
// @Tags favorites
// @Accept json
// @Produce json
// @Param work_id path string true "Work ID"
// @Success 200 {object} schema.IsFavoriteResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security BearerAuth
// @Router /auth/works/:work_id/favorite/is-favorite [get]
func (fc *FavoriteController) IsFavorite(c echo.Context) error {
	workIDStr := c.Param("work_id")
	workID, err := uuid.Parse(workIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work ID format")
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*schema.JWTCustomClaims)
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	isFavorite := fc.favoriteUsecase.IsFavorite(c.Request().Context(), workID, userID)
	return c.JSON(http.StatusOK, schema.IsFavoriteResponse{IsFavorite: isFavorite})
}

func handleFavoriteError(c echo.Context, err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
