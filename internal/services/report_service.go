package services

import (
	"context"
	"fmt"
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
	"strukit-services/pkg/report"

	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
)

func NewReport(receiptRepo *repository.ReceiptRepository, projectRepo *repository.ProjectRepository, reportManager *report.ReportManager) *ReportService {
	return &ReportService{
		receiptRepo:   receiptRepo,
		projectRepo:   projectRepo,
		reportManager: reportManager,
	}
}

type ReportService struct {
	receiptRepo   *repository.ReceiptRepository
	projectRepo   *repository.ProjectRepository
	reportManager *report.ReportManager
}

func (r *ReportService) DownloadExcelFile(ctx context.Context) (*excelize.File, error) {
	g, ctx := errgroup.WithContext(ctx)

	type result struct {
		project  *models.Project
		receipts []*models.Receipt
	}

	res := &result{}
	g.SetLimit(2)
	g.Go(func() error {
		project, err := r.projectRepo.GetProjectByID(ctx)
		if err != nil {
			return fmt.Errorf("[ReportService.DownloadExcelFile] error when GetProjectByID, err : %w", err)
		}

		res.project = project
		return nil
	})

	g.Go(func() error {
		receipts, err := r.receiptRepo.GetReceiptByProjectID(ctx)
		if err != nil {
			return fmt.Errorf("[ReportService.DownloadExcelFile] error when GetReceiptByProjectID, err : %w", err)
		}

		res.receipts = receipts
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	file, err := r.reportManager.GenerateExcel(res.project, res.receipts)
	if err != nil {
		return nil, fmt.Errorf("[ReportService.DownloadExcelFile] error when GenerateExcel, err : %w", err)
	}

	return file, nil
}
