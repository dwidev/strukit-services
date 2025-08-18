package services

import (
	"context"
	"net/http"
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
	"strukit-services/pkg/llm"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/responses"
)

func NewReceipt(llm *llm.Manager, repo *repository.ReceiptRepository, duplicateService *DuplicateDetectionService) *ReceiptService {
	return &ReceiptService{
		ReceiptRepository:         repo,
		Llm:                       llm,
		DuplicateDetectionService: duplicateService,
	}
}

type ReceiptService struct {
	*repository.ReceiptRepository
	Llm                       *llm.Manager
	DuplicateDetectionService *DuplicateDetectionService
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
	// 1. Extract data using AI
	r.Llm.Context = ctx
	receipt, err := r.Llm.ScanReceiptFromOCR(rawOcr)
	if err != nil {
		return nil, err
	}

	if !receipt.AIResponse.Success {
		return nil, responses.Err(http.StatusBadRequest, receipt.AIResponse.Message)
	}

	// 2. Convert to model
	model := receipt.Model()

	// 3. Check for duplicates

	// 4. Save to database
	res, err := r.ReceiptRepository.Save(ctx, model)
	if err != nil {
		logger.Log.Service(ctx, model).WithField("ai_response", receipt).WithField("error", err).Error("error when save receipt with success extraction data")
		return nil, err
	}

	return res, nil
}

func (r *ReceiptService) Scan(ctx context.Context, image []byte) (*models.Receipt, error) {
	// 1. Extract data using AI
	r.Llm.Context = ctx
	receipt, err := r.Llm.ScanReceiptWithImage(image)
	if err != nil {
		return nil, err
	}

	if !receipt.AIResponse.Success {
		return nil, responses.BodyErr(receipt.AIResponse.Message)
	}

	// 2. Convert to model
	model := receipt.Model()

	// 3. Check for duplicates
	err = r.CheckingDuplicates(ctx, model)
	if err != nil {
		return nil, err
	}

	// 4. Save to database
	res, err := r.ReceiptRepository.Save(ctx, model)
	if err != nil {
		logger.Log.Service(ctx, model).WithField("ai_response", receipt).WithField("error", err).Error("error when save receipt with success extraction data")
		return nil, err
	}

	return res, nil
}
