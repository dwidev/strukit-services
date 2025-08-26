package dto

type SpendingType string

const (
	Daily  SpendingType = "daily"
	Weekly SpendingType = "weekly"
)

type BudgetFilterRequest struct {
	Type SpendingType `json:"type" validate:"required"`
}
