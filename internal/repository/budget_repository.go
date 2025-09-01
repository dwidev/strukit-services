package repository

import (
	"context"
	"fmt"
	"math"
	"strukit-services/internal/dto"
	"strukit-services/internal/models"
	"strukit-services/pkg/budget"
	"strukit-services/pkg/helper"

	"golang.org/x/sync/errgroup"
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

func (b *BudgetRepository) GetBudgetByCategories(ctx context.Context) ([]dto.CategoryBudgetResponse, error) {
	q := `
		select 
			c.id,
			c."name",
			count(r.id) as total_receipt,
			round(coalesce(SUM(r.total_amount), 0), 2) as total_spent, 
			round(coalesce(avg(r.total_amount), 0), 2) as average_spent,
			round(coalesce(max(r.total_amount), 0), 2) as highest_transaction,
			round(coalesce(min(r.total_amount), 0), 2) as lowest_transaction
		from categories c
		left join receipts r on c.id = r.category_id
		and r.project_id = ?
		group by c."name", c.id
		order by total_spent desc
	`
	var category []dto.CategoryBudgetResponse
	if err := b.db.Raw(q, b.ProjectID(ctx)).Find(&category).Error; err != nil {
		return nil, fmt.Errorf("[BudgetRepository.GetBudgetSpending] error get daily spends, err : %s", err)
	}

	return category, nil
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

		if err := tx.Model(models.Receipt{}).Select("COALESCE(SUM(total_amount), 0)").Where("id = ? AND user_id = ?", projectId, userId).Scan(&totalSpent).Error; err != nil {
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
	remainingDays := int(project.EndDate.Sub(*b.DateNow()).Hours() / 24)

	projectStatusWording := b.getProjectBudgetStatusWording(&project, *totalSpent)

	budget := &dto.BudgetTrackingResponse{
		UserID:               userId,
		ProjectID:            projectId,
		BudgetAmount:         project.TotalBudget,
		TotalSpent:           *totalSpent,
		ProjectStatusWording: projectStatusWording,
		RemainingDays:        remainingDays,
		RemainingBudget:      remainingBudget,
		SpentPercentage:      spentPercentage,
		Receipts:             receipt,
	}

	return budget, nil
}

func (b *BudgetRepository) GetBudgetDetails(ctx context.Context, filter *dto.BudgetFilterRequest) (*dto.BudgetTrackingResponse, error) {
	type budgetResults struct {
		summary    *dto.BudgetTrackingResponse
		spending   []dto.BudgetSpendingData
		categories []dto.CategoryBudgetResponse
	}

	results := &budgetResults{}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		results.summary, err = b.GetBudgetSummary(ctx)
		if err != nil {
			return fmt.Errorf("[BudgetRepository.GetBudgetDetails] error when GetBudgetSummary: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		results.spending, err = b.GetBudgetSpending(ctx, filter)
		if err != nil {
			return fmt.Errorf("[BudgetRepository.GetBudgetDetails] error when GetBudgetSpending: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		results.categories, err = b.GetBudgetByCategories(ctx)
		if err != nil {
			return fmt.Errorf("[BudgetRepository.GetBudgetDetails] error when GetBudgetByCategories: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	budget := results.summary
	budget.Categories = results.categories
	budget.Spending = &dto.BudgetSpending{
		Type: filter.Type,
		Data: results.spending,
	}

	b.buildBurnRate(budget)
	b.buildBudgetProjection(filter, budget)

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

func (b *BudgetRepository) determineBudgetStatus(project *models.Project, totalSpent float64) budget.Status {
	usageRatio := math.Round((totalSpent/project.TotalBudget)*10) / 10

	if usageRatio >= 0.90 {
		return budget.Critical
	}

	if usageRatio > 0.75 {
		return budget.Warning
	}

	if usageRatio > 0.50 {
		return budget.Caution
	}

	return budget.OnTrack
}

// getProjectBudgetStatusWording returns overall project budget status wording
func (b *BudgetRepository) getProjectBudgetStatusWording(project *models.Project, totalSpent float64) dto.ProjectStatusWording {
	totalBudget := project.TotalBudget
	status := project.Status

	spentPercentage := (totalSpent / totalBudget) * 100
	remainingBudget := totalBudget - totalSpent

	budgetStatus := b.determineBudgetStatus(project, totalSpent)

	wording := dto.ProjectStatusWording{
		Status:       status,
		BudgetStatus: budgetStatus,
	}

	switch status {
	case models.ProjectStatusOver:
		overAmount := totalSpent - totalBudget
		wording.Title = "🚨 Budget Terlampaui"
		wording.Message = fmt.Sprintf("Total pengeluaran melebihi budget sebesar %s (%.1f%% dari budget)", helper.ParseToIDR(overAmount), spentPercentage)
		wording.ActionMessage = "Segera evaluasi pengeluaran dan pertimbangkan penambahan budget atau pengurangan scope"
		wording.Severity = "critical"
		wording.Color = "#D32F2F"

	case models.ProjectStatusActive:
		switch budgetStatus {
		case budget.Critical:
			wording.Title = "⚠️ Budget Hampir Habis"
			wording.Message = fmt.Sprintf("Sudah menggunakan %.1f%% budget. Tersisa %s", spentPercentage, helper.ParseToIDR(remainingBudget))
			wording.ActionMessage = "Kontrol ketat pengeluaran selanjutnya"
			wording.Severity = "high"
			wording.Color = "#FF5722"

		case budget.Warning:
			wording.Title = "🟡 Perhatian Budget"
			wording.Message = fmt.Sprintf("Sudah menggunakan %.1f%% budget. Tersisa %.0f",
				spentPercentage, remainingBudget)
			wording.ActionMessage = "Monitor pengeluaran dengan cermat"
			wording.Severity = "medium"
			wording.Color = "#FF9800"

		case budget.Caution:
			wording.Title = "📊 Budget Berjalan Normal"
			wording.Message = fmt.Sprintf("Sudah menggunakan %.1f%% budget. Tersisa %s", spentPercentage, helper.ParseToIDR(remainingBudget))
			wording.ActionMessage = "Pertahankan pola pengeluaran saat ini"
			wording.Severity = "low"
			wording.Color = "#4CAF50"

		default:
			wording.Title = "💚 Budget Aman"
			wording.Message = fmt.Sprintf("Baru menggunakan %d%% budget. Tersisa %s", int(spentPercentage), helper.ParseToIDR(remainingBudget))
			wording.ActionMessage = "Budget masih sangat aman untuk pengeluaran selanjutnya"
			wording.Severity = "info"
			wording.Color = "#2196F3"
		}

	case models.ProjectStatusCompleted:
		if spentPercentage <= 100 {
			wording.Title = "🎉 Project Selesai - Budget Terkontrol"
			wording.Message = fmt.Sprintf("Project selesai dengan menggunakan %.1f%% budget (%s)", spentPercentage, helper.ParseToIDR(totalSpent))
			wording.ActionMessage = "Sisa budget bisa digunakan untuk project lain"
		} else {
			overAmount := totalSpent - totalBudget
			wording.Title = "⚠️ Project Selesai - Over Budget"
			wording.Message = fmt.Sprintf("Project selesai dengan over budget %s (%.1f%% dari budget)", helper.ParseToIDR(overAmount), spentPercentage)
			wording.ActionMessage = "Evaluasi penyebab over budget untuk project selanjutnya"
		}
		wording.Severity = "info"
		wording.Color = "#9C27B0"

	case models.ProjectStatusArchived:
		wording.Title = "⏸️ Project Ditangguhkan"
		wording.Message = fmt.Sprintf("Project ditangguhkan dengan penggunaan budget %.1f%% (Rp %.0f)",
			spentPercentage, totalSpent)
		wording.ActionMessage = "Project dapat dilanjutkan kapan saja"
		wording.Severity = "info"
		wording.Color = "#607D8B"

	case models.ProjectStatusDeleted:
		wording.Title = "⛔️ Project Dihapus"
		wording.Message = fmt.Sprintf("Project dihapus dengan penggunaan budget %.1f%% (%s)", spentPercentage, helper.ParseToIDR(totalSpent))
		wording.ActionMessage = fmt.Sprintf("Project akan lenyap pada tanggal %s", project.DeletedAt.AddDate(0, 0, 15))
		wording.Severity = "info"
		wording.Color = "#FEF4444"

	default:
		wording.Title = "📊 Status Budget"
		wording.Message = "Status budget tidak diketahui"
		wording.ActionMessage = ""
		wording.Severity = "info"
		wording.Color = "#9E9E9E"
	}

	return wording
}

func (b *BudgetRepository) buildBurnRate(budget *dto.BudgetTrackingResponse) {
	if budget.Spending == nil || len(budget.Spending.Data) == 0 {
		budget.BurnRate = 0
		return
	}

	var allSpend float64
	for _, d := range budget.Spending.Data {
		allSpend += d.TotalAmount
	}

	burnRate := allSpend / float64(len(budget.Spending.Data))
	budget.BurnRate = burnRate
}

// method for create budget projections
func (b *BudgetRepository) buildBudgetProjection(filter *dto.BudgetFilterRequest, budget *dto.BudgetTrackingResponse) {
	if len(budget.Spending.Data) <= 1 {
		budget.Projections = nil
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
