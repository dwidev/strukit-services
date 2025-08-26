package dto

import (
	"strukit-services/pkg/budget"
	"time"

	"github.com/google/uuid"
)

type BudgetTrackingResponse struct {
	UserID          uuid.UUID     `json:"userId"`
	ProjectID       uuid.UUID     `json:"projectId"`
	BudgetAmount    float64       `json:"budgetAmount"`
	TotalSpent      float64       `json:"totalSpent"`
	DailyBurnRate   float64       `json:"dailyBurnRate"`
	RemainingDays   int           `json:"remainingDays"`
	RemainingBudget float64       `json:"remainingBudget"`
	SpentPercentage float64       `json:"spentPercentage"`
	BudgetStatus    budget.Status `json:"budgetStatus"`
	Receipts        BudgetReceipt `json:"total"`
	// AlertLevel      string        `json:"alertLevel"`
	// CategoryBreakdown []CategorySpending `json:"category_breakdown"`
	Spending    *BudgetSpending   `json:"spending,omitempty"`
	Projections *BudgetProjection `json:"projections,omitempty"`
	// Alert       []BudgetAlert `json:"alerts"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
}

type BudgetReceipt struct {
	Receipts int `json:"receipts"`
	Items    int `json:"items"`
}

type BudgetSpending struct {
	Type SpendingType          `json:"type"`
	Data *[]BudgetSpendingData `json:"data"`
}

type BudgetSpendingData struct {
	Date         time.Time `json:"date"`
	TotalAmount  float64   `json:"totalAmount"`
	Average      float64   `json:"average"`
	TotalReceipt float64   `json:"totalReceipt"`
}

type BudgetProjection struct {
	DailyBurnRate              float64   `json:"dailyBurnRate"`
	RemainingEstimedCompletion int       `json:"remainingEstimedCompletion"`
	EstimedCompletionDate      time.Time `json:"estimedCompletionDate"`
}
