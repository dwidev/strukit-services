package dto

type CategoryBudgetResponse struct {
	TotalAmount  float64 `json:"totalAmount"`
	CategoryName string  `json:"name"`
}
