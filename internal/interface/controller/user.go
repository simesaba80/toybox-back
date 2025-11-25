package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

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
