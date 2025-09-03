package controller

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type WorkController struct {
	workUsecase *usecase.WorkUseCase
}

func NewWorkController(workUsecase *usecase.WorkUseCase) *WorkController {
	return &WorkController{
		workUsecase: workUsecase,
	}
}

func (wc *WorkController) GetAllWorks(c echo.Context) error {
	works, err := wc.workUsecase.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, "Failed to retrieve works")
	}

	response := make([]schema.GetWorkOutput, len(works))
	for i, work := range works {
		response[i] = schema.ToWorkResponse(work)
	}

	return c.JSON(200, schema.WorkListResponse{Works: response})
}

func (wc *WorkController) GetWorkByID(c echo.Context) error {
	idStr := c.Param("work_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid work ID format")
	}

	work, err := wc.workUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(404, "Work not found")
		}
		return echo.NewHTTPError(500, "Failed to retrieve work details")
	}

	return c.JSON(200, schema.ToWorkResponse(work))
}

func (wc *WorkController) CreateWork(c echo.Context) error {
	var input schema.CreateWorkInput
	if err := c.Bind(&input); err != nil {
		c.Logger().Error("Bind error:", err)
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// リクエストボディからuser_idを取得し、UUIDにパース
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.Logger().Error("Invalid UserID format:", err)
		return echo.NewHTTPError(400, "Invalid UserID format")
	}

	createdWork, err := wc.workUsecase.CreateWork(
		c.Request().Context(),
		input.Title,
		input.Description,
		input.DescriptionHTML,
		input.Visibility,
		userID,
	)
	if err != nil {
		c.Logger().Error("WorkUseCase.CreateWork error:", err)
		return echo.NewHTTPError(500, "Failed to create work")
	}

	return c.JSON(201, schema.ToCreateWorkOutput(createdWork))
}
