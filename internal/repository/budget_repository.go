package repository

import (
	"context"
	"fmt"
	"math"
	"strukit-services/internal/dto"
	"strukit-services/internal/models"
	"strukit-services/pkg/budget"

	"gorm.io/gorm"
)

func NewBudget(base *BaseRepository) *BudgetRepository {
	return &BudgetRepository{
		BaseRepository: base,
	}
}

type BudgetRepository struct {
	*BaseRepository
}

func (b *BudgetRepository) GetBudgetSummary(ctx context.Context) (*dto.BudgetTrackingResponse, error) {
	projectId := b.ProjectID(ctx)
	userId := b.UserID(ctx)

	var project models.Project
	var totalSpent *float64
	var remainingBudget float64
	var spentPercentage float64
	var receipt dto.BudgetReceipt

	err := b.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Model(models.Project{}).Select("total_budget, status, end_date").Where("id = ? AND user_id = ?", projectId, userId).First(&project).Error; err != nil {
			return fmt.Errorf("[BudgetRepository.GetBudgetProjectDetails] error get total budget & status, err : %s", err)
		}

		if err := tx.Model(models.Receipt{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalSpent).Error; err != nil {
			return fmt.Errorf("[BudgetRepository.GetBudgetProjectDetails] error get total spend, err : %s", err)
		}

		q := `SELECT 
				(SELECT COUNT(*) FROM receipts r WHERE r.project_id = ? AND r.user_id = ?) AS receipts,
				(SELECT COUNT(*) FROM receipt_items ri JOIN receipts r ON r.id = ri.receipt_id WHERE r.project_id = ? AND r.user_id = ?) AS items`
		if err := tx.Raw(q, projectId, userId, projectId, userId).Scan(&receipt).Error; err != nil {
			return fmt.Errorf("[BudgetRepository.GetBudgetProjectDetails] error get total budget & status, err : %s", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	remainingBudget = math.Round(project.TotalBudget - *totalSpent)
	spentPercentage = math.Round((*totalSpent / project.TotalBudget) * 100)
	bs := b.determineBudgetStatus(&project, int(spentPercentage))
	remainingDays := int(project.EndDate.Sub(*b.DateNow()).Hours() / 24)

	budget := &dto.BudgetTrackingResponse{
		UserID:          userId,
		ProjectID:       projectId,
		BudgetAmount:    project.TotalBudget,
		TotalSpent:      *totalSpent,
		RemainingDays:   remainingDays,
		RemainingBudget: remainingBudget,
		SpentPercentage: spentPercentage,
		BudgetStatus:    bs,
		Receipts:        receipt,
	}

	return budget, nil
}

func (b *BudgetRepository) GetBudgetDetails(ctx context.Context, filter *dto.BudgetFilterRequest) (*dto.BudgetTrackingResponse, error) {
	summary, err := b.GetBudgetSummary(ctx)
	if err != nil {
		return nil, err
	}

	spending, err := b.GetBudgetSpending(ctx, filter)
	if err != nil {
		return nil, err
	}

	summary.Spending = &dto.BudgetSpending{
		Type: filter.Type,
		Data: &spending,
	}
	budget := summary

	b.buildBurnRate(summary)
	b.buildBudgetProjection(filter, summary)

	return budget, nil
}

func (b *BudgetRepository) GetBudgetSpending(ctx context.Context, filter ...*dto.BudgetFilterRequest) ([]dto.BudgetSpendingData, error) {
	request := &dto.BudgetFilterRequest{
		Type: dto.Daily,
	}

	if len(filter) > 0 {
		request = filter[0]
	}

	q := `
		select 
			date_trunc(?, r.created_at) as date,
			round(coalesce(SUM(r.total_amount), 0), 2) as total_amount,
			count(r.id) total_receipt,
			round(avg(total_amount), 2) as average
		from receipts r 
		where r.project_id = ? AND r.user_id = ?
		group by date 
		order by date 
	`
	var spending []dto.BudgetSpendingData
	if err := b.db.Raw(q, request.Filter(), b.ProjectID(ctx), b.UserID(ctx)).Find(&spending).Error; err != nil {
		return nil, fmt.Errorf("[BudgetRepository.GetBudgetSpending] error get daily spends, err : %s", err)
	}

	return spending, nil
}

func (b *BudgetRepository) determineBudgetStatus(project *models.Project, spentPercentage int) budget.Status {
	if spentPercentage > 100 {
		return budget.OverBudget
	}

	if b.Now().After(*project.EndDate) {
		return budget.Completed
	}

	return budget.Healthty
}

func (b *BudgetRepository) buildBurnRate(budget *dto.BudgetTrackingResponse) {
	if budget.Spending == nil {
		return
	}

	var allSpend float64
	for _, d := range *budget.Spending.Data {
		allSpend += d.TotalAmount
	}

	burnRate := allSpend / float64(len(*budget.Spending.Data))
	budget.BurnRate = burnRate
}

// method for create budget projections
func (b *BudgetRepository) buildBudgetProjection(filter *dto.BudgetFilterRequest, budget *dto.BudgetTrackingResponse) {
	if budget.Spending.Data == nil || len(*budget.Spending.Data) < 2 {
		return
	}

	burnRate := budget.BurnRate
	remaining := int(budget.RemainingBudget / burnRate)

	estimed := b.DateNow().AddDate(0, 0, remaining)

	if filter.Weekly() {
		estimed = b.DateNow().AddDate(0, 0, 7*remaining)
	}

	if filter.Yearly() {
		estimed = b.DateNow().AddDate(0, 0, 365*remaining)
	}

	budget.Projections = &dto.BudgetProjection{
		Type:                       filter.Type,
		BurnRate:                   burnRate,
		RemainingEstimedCompletion: remaining,
		EstimedCompletionDate:      estimed,
	}
}
