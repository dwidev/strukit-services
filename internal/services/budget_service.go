package services

import (
	"context"
	"strukit-services/internal/dto"
	"strukit-services/internal/repository"
)

func NewBudget(budgetRepo *repository.BudgetRepository) *BudgetService {
	return &BudgetService{BudgetRepository: budgetRepo}
}

type BudgetService struct {
	*repository.BudgetRepository
}

func (b *BudgetService) GetBudgetSummary(ctx context.Context) (*dto.BudgetTrackingResponse, error) {
	response, err := b.BudgetRepository.GetBudgetSummary(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (b *BudgetService) GetBudgetDetails(ctx context.Context, filter *dto.BudgetFilterRequest) (*dto.BudgetTrackingResponse, error) {
	response, err := b.BudgetRepository.GetBudgetDetails(ctx, filter)
	if err != nil {
		return nil, err
	}

	return response, nil
}
