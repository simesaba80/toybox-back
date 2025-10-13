package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/simesaba80/toybox-back/internal/interface/schema"
	"github.com/simesaba80/toybox-back/internal/usecase"
)

type CommentController struct {
	commentUsecase usecase.CommentUsecase
}

func NewCommentController(commentUsecase usecase.CommentUsecase) *CommentController {
	return &CommentController{
		commentUsecase: commentUsecase,
	}
}

// GetCommentsByWorkID godoc
// @Summary Get comments for a work
// @Description Get all comments for a specific work
// @Tags comments
// @Produce json
// @Param work_id path string true "Work ID"
// @Success 200 {array} schema.CommentResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /works/{work_id}/comments [get]
func (cc *CommentController) GetCommentsByWorkID(c echo.Context) error {
	workIDStr := c.Param("work_id")
	workID, err := uuid.Parse(workIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid work ID format")
	}

	comments, err := cc.commentUsecase.GetCommentsByWorkID(c.Request().Context(), workID)
	if err != nil {
		c.Logger().Error("CommentUsecase.GetCommentsByWorkID error:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve comments")
	}

	return c.JSON(http.StatusOK, schema.ToCommentListResponse(comments))
}
