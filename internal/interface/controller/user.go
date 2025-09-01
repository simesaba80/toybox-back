package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/presenter"
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

func (uc *UserController) CreateUser(c echo.Context) error {
	var user presenter.CreateUserInput
	if err := c.Bind(&user); err != nil {
		return err
	}

	createdUser, err := uc.userusecase.CreateUser(c.Request().Context(), user.Name, user.Email, user.PasswordHash, user.DisplayName, user.AvatarURL)
	if err != nil {
		return err
	}

	return c.JSON(201, presenter.ToUserResponse(createdUser))
}

func (uc *UserController) GetAllUsers(c echo.Context) error {
	users, err := uc.userusecase.GetAllUser(c.Request().Context())
	if err != nil {
		return err
	}

	response := make([]presenter.GetUserOutput, len(users))
	for i, user := range users {
		response[i] = presenter.ToUserResponse(user)
	}

	return c.JSON(200, presenter.UserListResponse{Users: response})
}
