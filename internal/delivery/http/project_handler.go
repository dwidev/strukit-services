package http

import (
	"net/http"
	"strukit-services/internal/dto"
	"strukit-services/internal/services"
	"strukit-services/pkg/responses"

	"github.com/gin-gonic/gin"
)

func NewProject(
	base *BaseHandler,
	projectService *services.ProjectService,
	budgetService *services.BudgetService,
) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler:    base,
		ProjectService: projectService,
		BudgetService:  budgetService,
	}
}

type ProjectHandler struct {
	*BaseHandler
	*services.ProjectService
	*services.BudgetService
}

func (a *ProjectHandler) GetBudgetSummary(c *gin.Context) {
	ctx := c.Request.Context()
	results, err := a.BudgetService.GetBudgetSummary(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, responses.DataResponse{
		StatusCode: http.StatusOK,
		Data:       results,
	})
}

func (a *ProjectHandler) GetBudgetDetails(c *gin.Context) {
	ctx := c.Request.Context()
	var body dto.BudgetFilterRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if msg := a.AppValidator.Valid(&body); len(msg) > 0 {
		err := responses.BodyErr(msg)
		c.Error(err)
		return
	}
	results, err := a.BudgetService.GetBudgetDetails(ctx, &body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, responses.DataResponse{
		StatusCode: http.StatusOK,
		Data:       results,
	})
}

func (a *ProjectHandler) GetProjectByID(c *gin.Context) {
	ctx := c.Request.Context()
	results, err := a.ProjectService.GetProjectByID(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, responses.DataResponse{
		StatusCode: http.StatusOK,
		Data:       results,
	})
}

func (a *ProjectHandler) All(c *gin.Context) {
	ctx := c.Request.Context()

	results, err := a.ProjectService.All(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, responses.DataResponse{
		StatusCode: http.StatusOK,
		Data:       results,
	})
}

func (a *ProjectHandler) SoftDelete(c *gin.Context) {
	projectID := c.Param("id")
	ctx := c.Request.Context()

	if projectID == "" {
		err := responses.BodyErr("project id cannot be empty")
		c.Error(err)
		return
	}

	err := a.ProjectService.SoftDelete(ctx, projectID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projectId": projectID,
		"message":   "Project deleted successfully",
	})
}

func (a *ProjectHandler) CreateNewProject(c *gin.Context) {
	ctx := c.Request.Context()
	var body dto.CreateProjectDto

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Error(err)
		return
	}

	if msg := a.AppValidator.Valid(&body); len(msg) > 0 {
		err := responses.BodyErr(msg)
		c.Error(err)
		return
	}

	response, err := a.ProjectService.CreateNewProject(ctx, &body)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
