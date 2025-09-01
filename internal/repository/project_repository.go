package repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"strukit-services/internal/models"
	appContext "strukit-services/pkg/context"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/responses"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewProject(base *BaseRepository) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: base,
	}
}

type ProjectRepository struct {
	*BaseRepository
}

func (p *ProjectRepository) CheckExistProject(ctx context.Context) error {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	projectId := ctx.Value(appContext.ProjectID).(uuid.UUID)

	var exist bool
	q := "SELECT EXISTS (SELECT 1 FROM projects WHERE id = ? AND user_id = ?)"
	if err := p.db.Raw(q, projectId, userId).Scan(&exist).Error; err != nil {
		return err
	}

	if !exist {
		return responses.Err(http.StatusNotFound, fmt.Sprintf("Project with id %s and user_id %s do not exist", projectId, userId))
	}

	return nil
}

func (p *ProjectRepository) GetProjectByID(ctx context.Context) (result *models.Project, err error) {
	projectId := ctx.Value(appContext.ProjectID).(uuid.UUID)
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	var project *models.Project
	res := p.db.First(&project, "id = ? AND user_id = ?", projectId, userId)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, responses.Err(http.StatusNotFound, fmt.Sprintf("not found project by project id %s", projectId))
		}
		logger.Log.DB(ctx).Errorf("error when get project by id : %s", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, responses.AppForbidden
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
		Updates(&models.Project{IsSoftDelete: true, DeletedAt: p.Now()})

	if result.Error != nil {
		return err
	}

	if result.RowsAffected == 0 {
		logger.Log.DB(ctx).Warnf("delete project error, not match user_id:%s and project_id:%s", userId, projectID)
		return responses.AppForbidden
	}

	return nil
}

func (p *ProjectRepository) CreateNewProject(ctx context.Context, project *models.Project) (res *models.Project, err error) {
	userId := ctx.Value(appContext.UserIDKey).(uuid.UUID)
	status := models.ProjectStatusActive
	project.UserID = userId
	project.Status = status
	if err = p.db.Create(project).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return nil, responses.Err(http.StatusConflict, fmt.Sprintf("Project with name %s already exist, please create different project name", project.Name))
		}

		logger.Log.DB(ctx).Errorf("error create project : %s", err)
		return nil, err
	}

	return project, nil
}
