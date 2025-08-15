package services

import (
	"context"
	"strukit-services/internal/dto"
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
)

func NewProject(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{ProjectRepository: projectRepo}
}

type ProjectService struct {
	*repository.ProjectRepository
}

func (a *ProjectService) GetProjectByID(ctx context.Context, projectID string) (results *models.Project, err error) {
	results, err = a.ProjectRepository.GetProjectByID(ctx, projectID)
	if err != nil {
		return
	}

	return
}

func (a *ProjectService) All(ctx context.Context) (results []*models.Project, err error) {
	results, err = a.ProjectRepository.All(ctx)
	if err != nil {
		return
	}

	return
}

func (a *ProjectService) SoftDelete(ctx context.Context, projectID string) (err error) {
	if err = a.ProjectRepository.SoftDelete(ctx, projectID); err != nil {
		return
	}

	return nil
}

func (a *ProjectService) CreateNewProject(ctx context.Context, dto *dto.CreateProjectDto) (token *models.Project, err error) {
	newProject := dto.Model()
	project, err := a.ProjectRepository.CreateNewProject(ctx, newProject)
	if err != nil {
		return nil, err
	}

	return project, nil
}
