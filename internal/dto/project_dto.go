package dto

import (
	"strukit-services/internal/models"
	"strukit-services/pkg/helper"

	"github.com/google/uuid"
)

type CreateProjectDto struct {
	Name        string `validate:"required" json:"Name"`
	Description string `validate:"required" json:"description"`
	TotalBudget int    `validate:"required" json:"totalBudget"`
	StartDate   string `validate:"required" json:"startDate"`
	EndDate     string `validate:"required" json:"endDate"`
}

func (c *CreateProjectDto) Model(userID uuid.UUID) *models.Project {
	return &models.Project{
		UserID:      userID,
		Name:        c.Name,
		Description: c.Description,
		TotalBudget: c.TotalBudget,
		StartDate:   helper.ParseToDate(c.StartDate),
		EndDate:     helper.ParseToDate(c.EndDate),
	}
}
