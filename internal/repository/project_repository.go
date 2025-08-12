package repository

import (
	"context"
	"strukit-services/internal/models"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/responses"

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

func (p *ProjectRepository) All(ctx context.Context) (results []*models.Project, err error) {
	userId := ctx.Value(appContext.UserIDKey).(string)
	var projects []*models.Project
	if err = p.db.Find(&projects, "user_id = ?", uuid.MustParse(userId)).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectRepository) SoftDelete(ctx context.Context, projectID string) (err error) {
	userId := ctx.Value(appContext.UserIDKey).(string)
	result := p.db.Model(&models.Project{}).
		Where("id = ? AND user_id = ?", uuid.MustParse(projectID), uuid.MustParse(userId)).
		Update("is_soft_delete", true)

	if result.Error != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return responses.Forbidden()
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
