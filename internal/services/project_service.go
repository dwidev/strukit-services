package services

import (
	"context"
	"strukit-services/internal/dto"
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
	appContext "strukit-services/pkg/context"

	"github.com/google/uuid"
)

func NewProject(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepository: projectRepo}
}

type ProjectService struct {
	*repository.ProjectRepository
}

func (a *ProjectService) CreateNewProject(ctx context.Context, dto *dto.CreateProjectDto) (token *models.Project, err error) {
	userId := ctx.Value(appContext.UserIDKey).(string)
	newProject := dto.Model(uuid.MustParse(userId))
	project, err := a.ProjectRepository.CreateNewProject(newProject)
	if err != nil {
		return nil, err
	}

	return project, nil
}
