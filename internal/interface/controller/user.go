package controller

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	domainerrors "github.com/simesaba80/toybox-back/internal/domain/errors"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type UserController struct {
	userusecase usecase.IUserUseCase
}

func NewUserController(userusecase usecase.IUserUseCase) *UserController {
	return &UserController{
		userusecase: userusecase,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce  json
// @Success 200 {object} schema.UserListResponse
// @Router /users [get]
func (uc *UserController) GetAllUsers(c echo.Context) error {
	users, err := uc.userusecase.GetAllUser(c.Request().Context())
	if err != nil {
		return err
	}

	response := make([]schema.GetUserOutput, len(users))
	for i, user := range users {
		response[i] = schema.ToUserResponse(user)
	}

	return c.JSON(http.StatusOK, schema.UserListResponse{Users: response})
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} schema.GetUserOutput
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [get]
func (uc *UserController) GetUserByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なリクエストです")
	}
	user, err := uc.userusecase.GetByUserID(c.Request().Context(), id)
	if err != nil {
		return handleUserError(err)
	}
	return c.JSON(http.StatusOK, schema.ToUserResponse(user))
}

func handleUserError(err error) error {
	switch {
	case errors.Is(err, domainerrors.ErrUserNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
}
