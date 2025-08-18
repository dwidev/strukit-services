package services

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
	"strukit-services/pkg/llm"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/responses"
)

func NewReceipt(llm *llm.Manager, repo *repository.ReceiptRepository, projectRepo *repository.ProjectRepository, duplicateService *DuplicateDetectionService) *ReceiptService {
	return &ReceiptService{
		ProjectRepository:         projectRepo,
		ReceiptRepository:         repo,
		Llm:                       llm,
		DuplicateDetectionService: duplicateService,
	}
}

type ReceiptService struct {
	*repository.ProjectRepository
	*repository.ReceiptRepository
	Llm                       *llm.Manager
	DuplicateDetectionService *DuplicateDetectionService
}

func (r *ReceiptService) GetReceiptByProjectID(ctx context.Context) (receipts []*models.Receipt, err error) {
	if receipts, err = r.ReceiptRepository.GetReceiptByProjectID(ctx); err != nil {
		return nil, fmt.Errorf("[main.(*ReceiptService).GetReceiptByProjectID] error when get receipt %w", err)
	}

	return receipts, nil
}

func (r *ReceiptService) CheckingDuplicates(ctx context.Context, receipt *models.Receipt) error {
	checked, err := r.DuplicateDetectionService.Checking(ctx, receipt)
	if err != nil {
		logger.Log.Service(ctx, receipt).WithField("error", err).Error("error checking for duplicates")
		return err
	}

	if checked.HasDuplicates {
		switch checked.RecommendedAction {
		case Block:
			logger.Log.Service(ctx, receipt).Errorf("Detection duplicated receipt")
			return responses.Err(http.StatusConflict, "Detection duplicated receipt")
		case Warning:
			logger.Log.Service(ctx, receipt).Errorf("Detection duplicated receipt please correction before save")
		case Review:
			logger.Log.Service(ctx, receipt).Errorf("Not sure detection duplicated receipt please review before save")
		}
	}

	return nil
}

func (r *ReceiptService) ScanFromOCR(ctx context.Context, rawOcr string) (*models.Receipt, error) {
	model, err := r.processingData(ctx, func() (*llm.ReceiptResponse, error) {
		return r.Llm.ScanReceiptFromOCR(rawOcr)
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ReceiptService) Scan(ctx context.Context, image []byte) (*models.Receipt, error) {
	model, err := r.processingData(ctx, func() (*llm.ReceiptResponse, error) {
		return r.Llm.ScanReceiptWithImage(image)
	})

	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *ReceiptService) processingData(ctx context.Context, fn func() (*llm.ReceiptResponse, error)) (*models.Receipt, error) {
	// 1. Checking project exists
	err := r.ProjectRepository.CheckExistProject(ctx)
	if err != nil {
		fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		return nil, fmt.Errorf("error when checking project exist at [main.(*ReceiptService).processingData] with %s Err : %w", fnName, err)
	}

	// 2. Generate LLM extraction
	r.Llm.Context = ctx
	receipt, err := fn()
	if err != nil {
		return nil, err
	}

	if !receipt.AIResponse.Success {
		return nil, responses.BodyErr(receipt.AIResponse.Message)
	}

	// 3. Checking duplicated data
	model := receipt.Model()
	err = r.CheckingDuplicates(ctx, model)
	if err != nil {
		logger.Log.Service(ctx, model).WithField("error", err).Error("error when check duplicated data")
		return nil, err
	}

	// 4. Saving receipt
	res, err := r.ReceiptRepository.Save(ctx, model)
	if err != nil {
		logger.Log.Service(ctx, model).WithField("error", err).Error("error when save receipt with success extraction data")
		return nil, err
	}

	return res, nil
}

func (r *ReceiptService) Delete(ctx context.Context) error {
	if err := r.ReceiptRepository.Delete(ctx); err != nil {
		return fmt.Errorf("[main.(*ReceiptService).Delete] error when delete receipt %w", err)
	}

	return nil
}
