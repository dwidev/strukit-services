package dto

import (
	"strukit-services/internal/models"
	"strukit-services/pkg/helper"
)

type CreateProjectDto struct {
	Name        string  `validate:"required" json:"Name"`
	Description string  `validate:"required" json:"description"`
	TotalBudget float64 `validate:"required" json:"totalBudget"`
	StartDate   string  `validate:"required" json:"startDate"`
	EndDate     string  `validate:"required" json:"endDate"`
}

func (c *CreateProjectDto) Model() *models.Project {
	return &models.Project{
		Name:        c.Name,
		Description: c.Description,
		TotalBudget: c.TotalBudget,
		StartDate:   helper.ParseToDate(&c.StartDate),
		EndDate:     helper.ParseToDate(&c.EndDate),
	}
}
