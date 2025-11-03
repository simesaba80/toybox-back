package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type WorkController struct {
	workUsecase usecase.IWorkUseCase
}

func NewWorkController(workUsecase usecase.IWorkUseCase) *WorkController {
	return &WorkController{
		workUsecase: workUsecase,
	}
}

// GetAllWorks godoc
// @Summary Get all works
// @Description Get all works with pagination
// @Tags works
// @Produce json
// @Param limit query int false "Limit per page (default: 20, max: 100)"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} schema.WorkListResponse
// @Router /works [get]
func (wc *WorkController) GetAllWorks(c echo.Context) error {
	var query schema.GetWorksQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters")
	}
	if err := c.Validate(&query); err != nil {
		return err
	}

	works, total, limit, page, err := wc.workUsecase.GetAll(c.Request().Context(), query.Limit, query.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve works")
	}

	response := make([]schema.GetWorkOutput, len(works))
	for i, work := range works {
		response[i] = schema.ToWorkResponse(work)
	}

	return c.JSON(http.StatusOK, schema.WorkListResponse{
		Works:      response,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
	})
}

// GetWorkByID godoc
// @Summary Get a work by ID
// @Description Get a work by ID
// @Tags works
// @Produce json
// @Param work_id path string true "Work ID"
// @Success 200 {object} schema.GetWorkOutput
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /works/{work_id} [get]
func (wc *WorkController) GetWorkByID(c echo.Context) error {
	idStr := c.Param("work_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work ID format")
	}

	work, err := wc.workUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Work not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve work details")
	}

	return c.JSON(http.StatusOK, schema.ToWorkResponse(work))
}

// CreateWork godoc
// @Summary Create a new work
// @Description Create a new work with the input payload
// @Tags works
// @Accept json
// @Produce json
// @Param work body schema.CreateWorkInput true "Work to create"
// @Success 201 {object} schema.CreateWorkOutput
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /works [post]
func (wc *WorkController) CreateWork(c echo.Context) error {
	var input schema.CreateWorkInput
	if err := c.Bind(&input); err != nil {
		c.Logger().Error("Bind error:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
    if err := c.Validate(&input); err != nil {
      return err
    }
	// リクエストボディからuser_idを取得し、UUIDにパース
	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.Logger().Error("Invalid UserID format:", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid UserID format")
	}

	createdWork, err := wc.workUsecase.CreateWork(
		c.Request().Context(),
		input.Title,
		input.Description,
		input.Visibility,
		userID,
	)
	if err != nil {
		c.Logger().Error("WorkUseCase.CreateWork error:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create work")
	}

	return c.JSON(http.StatusCreated, schema.ToCreateWorkOutput(createdWork))
}
