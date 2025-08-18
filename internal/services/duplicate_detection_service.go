package services

import (
	"context"
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
	"strukit-services/pkg/hash"
	"strukit-services/pkg/logger"
)

type DetectionAction string

var (
	Block   DetectionAction = "block"
	Warning DetectionAction = "warn"
	Review  DetectionAction = "review"
)

type DuplicateResult struct {
	HasDuplicates        bool
	RecommendedAction    DetectionAction
	DetecetionConfidence float64
	DuplicatedReceipts   []*models.Receipt
}

type DuplicateDetectionService struct {
	receiptRepo *repository.ReceiptRepository
}

func NewDuplicateDetectionService(receiptRepo *repository.ReceiptRepository) *DuplicateDetectionService {
	service := &DuplicateDetectionService{
		receiptRepo: receiptRepo,
	}

	return service
}

func (s *DuplicateDetectionService) Checking(ctx context.Context, receipt *models.Receipt) (res *DuplicateResult, err error) {
	s.hashReceipt(receipt)

	duplicates, err := s.receiptRepo.FindByFingerprint(ctx, receipt.Fingerprint)
	if err != nil {
		logger.Log.Service(ctx, receipt).WithField("error", err).Error("error checking for duplicates fingerprint")
		return nil, err
	}

	if len(duplicates) > 0 {
		logger.Log.Service(ctx, receipt).WithField("duplicates", duplicates).Info("duplicate receipt detected fingerprint")
		return s.generateResponse(duplicates, "fingerprint"), nil
	}

	duplicates, err = s.receiptRepo.FindByContentHash(ctx, receipt.ContentHash)
	if err != nil {
		logger.Log.Service(ctx, receipt).WithField("error", err).Error("error checking for duplicates content hash")
		return nil, err
	}

	if len(duplicates) > 0 {
		logger.Log.Service(ctx, receipt).WithField("duplicates", duplicates).Info("duplicate receipt detected content hash")
		return s.generateResponse(duplicates, "content_hash"), nil
	}

	duplicates, err = s.receiptRepo.FindSimilarReceipts(ctx, receipt)
	if err != nil {
		logger.Log.Service(ctx, receipt).WithField("error", err).Error("error checking for duplicates similarity")
		return nil, err
	}

	if len(duplicates) > 0 {
		logger.Log.Service(ctx, receipt).WithField("similarity", duplicates).Info("duplicate receipt detected similarity")
		return s.generateResponse(duplicates, "similarity"), nil
	}

	return s.generateResponse([]*models.Receipt{}), nil
}

func (s *DuplicateDetectionService) hashReceipt(receipt *models.Receipt) {
	receiptHasData := hash.ReceiptHashData{
		ProjectID:       receipt.ProjectID,
		MerchantName:    receipt.MerchantName,
		ReceiptNumber:   receipt.ReceiptNumber,
		TotalAmount:     receipt.SubTotal,
		TransactionDate: receipt.TransactionDate.Format("2006-01-02"),
		TransactionTime: receipt.TransactionTime.Format(),
		Items:           receipt.Items,
	}

	receipt.ContentHash = hash.GenerateContentHash(receiptHasData)
	receipt.Fingerprint = hash.GenerateFingerprint(receiptHasData.ProjectID, receipt.ContentHash)
}

func (s *DuplicateDetectionService) generateResponse(duplicates []*models.Receipt, f ...string) *DuplicateResult {
	from := "none"
	if len(f) > 0 {
		from = f[0]
	}

	if from == "fingerprint" {
		return &DuplicateResult{
			HasDuplicates:        true,
			DetecetionConfidence: 1.0,
			RecommendedAction:    "block",
			DuplicatedReceipts:   duplicates,
		}
	}

	if from == "content_hash" {
		return &DuplicateResult{
			HasDuplicates:        true,
			DetecetionConfidence: 0.80,
			RecommendedAction:    "warn_and_recheck",
			DuplicatedReceipts:   duplicates,
		}
	}

	if from == "similarity" {
		return &DuplicateResult{
			HasDuplicates:        true,
			DetecetionConfidence: 0.50,
			RecommendedAction:    "recheck",
			DuplicatedReceipts:   duplicates,
		}
	}

	return &DuplicateResult{
		HasDuplicates:        false,
		DetecetionConfidence: 0.0,
		DuplicatedReceipts:   nil,
	}
}
