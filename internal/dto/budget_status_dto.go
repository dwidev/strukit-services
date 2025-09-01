package dto

import (
	"strukit-services/internal/models"
	"strukit-services/pkg/budget"
)

// ProjectStatusWording for overall project status
type ProjectStatusWording struct {
	Status        models.ProjectStatus `json:"status"`
	BudgetStatus  budget.Status        `json:"budgetStatus"`
	Title         string               `json:"title"`
	Message       string               `json:"message"`
	ActionMessage string               `json:"actionMessage"`
	Severity      string               `json:"severity"`
	Color         string               `json:"color"`
}
