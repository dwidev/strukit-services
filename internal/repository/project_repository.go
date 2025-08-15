package repository

import (
	"context"
	"strukit-services/internal/models"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/logger"
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

func (p *ProjectRepository) GetProjectByID(ctx context.Context, projectId string) (result *models.Project, err error) {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	var project *models.Project
	res := p.db.First(&project, "id = ? AND user_id = ?", uuid.MustParse(projectId), userId)

	if res.Error != nil {
		logger.Log.Errorf("error when get project by id : %s", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, responses.Forbidden()
	}

	return project, nil
}

func (p *ProjectRepository) All(ctx context.Context) (results []*models.Project, err error) {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	var projects []*models.Project
	if err = p.db.Find(&projects, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *ProjectRepository) SoftDelete(ctx context.Context, projectID string) (err error) {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	result := p.db.Model(&models.Project{}).
		Where("id = ? AND user_id = ?", uuid.MustParse(projectID), userId).
		Update("is_soft_delete", true)

	if result.Error != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return responses.Forbidden()
	}

	return nil
}

func (p *ProjectRepository) CreateNewProject(ctx context.Context, project *models.Project) (res *models.Project, err error) {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	status := models.ProjectStatusActive
	project.UserID = userId
	project.Status = &status
	if err = p.db.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}
