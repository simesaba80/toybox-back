package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type UserController struct {
	userusecase *usecase.UserUseCase
}

func NewUserController(userusecase *usecase.UserUseCase) *UserController {
	return &UserController{
		userusecase: userusecase,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user body schema.CreateUserInput true "User to create"
// @Success 201 {object} schema.GetUserOutput
// @Router /users [post]
func (uc *UserController) CreateUser(c echo.Context) error {
	var user schema.CreateUserInput
	if err := c.Bind(&user); err != nil {
		return err
	}
    if err := c.Validate(&user); err != nil {
      return err
    }

	createdUser, err := uc.userusecase.CreateUser(c.Request().Context(), user.Name, user.Email, user.PasswordHash, user.DisplayName, user.AvatarURL)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, schema.ToUserResponse(createdUser))
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
