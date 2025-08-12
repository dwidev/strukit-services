package http

import (
	"strukit-services/internal/dto"
	"strukit-services/internal/services"
	"strukit-services/pkg/responses"

	"github.com/gin-gonic/gin"
)

func NewProject(
	base *BaseHandler,
	projectService *services.ProjectService,
) *ProjectHandler {
	return &ProjectHandler{
		BaseHandler:    base,
		ProjectService: projectService,
	}
}

type ProjectHandler struct {
	*BaseHandler
	*services.ProjectService
}

func (a *ProjectHandler) SoftDelete(c *gin.Context) {
	projectID := c.Param("id")

	if projectID == "" {
		err := responses.BodyErr("project id cannot be empty")
		c.Error(err)
		return
	}

	err := a.ProjectService.SoftDelete(projectID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
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

	c.JSON(200, response)
}
