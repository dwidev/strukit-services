package dto

import (
	"strukit-services/internal/models"
	"time"
)

type CreateProjectDto struct {
	Name        string     `validate:"required,min=5,max=155" json:"Name"`
	Description string     `validate:"required,max=200" json:"description"`
	TotalBudget float64    `validate:"required,max=999999999" json:"totalBudget"`
	StartDate   *time.Time `validate:"required" json:"startDate"`
	EndDate     *time.Time `validate:"required,gtfield=StartDate" json:"endDate"`
}

func (c *CreateProjectDto) Model() *models.Project {
	return &models.Project{
		Name:        c.Name,
		Description: c.Description,
		TotalBudget: c.TotalBudget,
		StartDate:   c.StartDate,
		EndDate:     c.EndDate,
	}
}
