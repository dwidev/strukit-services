package dto

type SpendingType string

const (
	Daily  SpendingType = "daily"
	Weekly SpendingType = "weekly"
	Yearly SpendingType = "yearly"
)

type BudgetFilterRequest struct {
	Type SpendingType `json:"type" validate:"required"`
}

func (b BudgetFilterRequest) Filter() string {
	if b.Type == Yearly {
		return "year"
	}

	if b.Type == Weekly {
		return "week"
	}

	return "day"
}

func (b BudgetFilterRequest) Daily() bool {
	return b.Type == Daily
}

func (b BudgetFilterRequest) Weekly() bool {
	return b.Type == Weekly
}

func (b BudgetFilterRequest) Yearly() bool {
	return b.Type == Yearly
}
