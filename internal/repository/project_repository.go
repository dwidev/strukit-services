package repository

import (
	"strukit-services/internal/models"

	"github.com/google/uuid"
)

func NewProject(base *BaseRepository) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: base,
	}
}

type ProjectRepository struct {
	*BaseRepository
}

func (p *ProjectRepository) SoftDelete(projectID string) (err error) {
	if err = p.db.Model(&models.Project{}).Where("id = ?", uuid.MustParse(projectID)).Update("is_soft_delete", true).Error; err != nil {
		return
	}

	return nil
}

func (p *ProjectRepository) CreateNewProject(project *models.Project) (res *models.Project, err error) {
	status := models.ProjectStatusActive
	project.Status = &status
	if err = p.db.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}
