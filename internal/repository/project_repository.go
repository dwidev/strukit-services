package repository

import "strukit-services/internal/models"

func NewProject(base *BaseRepository) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: base,
	}
}

type ProjectRepository struct {
	*BaseRepository
}

func (p *ProjectRepository) CreateNewProject(project *models.Project) (res *models.Project, err error) {
	status := models.ProjectStatusActive
	project.Status = &status
	if err = p.db.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}
