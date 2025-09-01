package dto

type CategoryBudgetResponse struct {
	Id                 string  `json:"id"`
	Name               string  `json:"name"`
	TotalReceipt       int     `json:"totalReceipt"`
	TotalSpent         float64 `json:"totalSpent"`
	AverageSpent       float64 `json:"averageSpent"`
	HighestTransaction float64 `json:"highestTransaction"`
	LowestTransaction  float64 `json:"lowestTransaction"`
}
