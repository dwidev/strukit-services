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

func NewReceipt(llm *llm.Manager, repo *repository.ReceiptRepository) *ReceiptService {
	return &ReceiptService{
		ReceiptRepository: repo,
		Llm:               llm,
	}
}

type ReceiptService struct {
	*repository.ReceiptRepository
	Llm *llm.Manager
}

func (r *ReceiptService) ScanFromOCR(ctx context.Context, rawOcr string) (*models.Receipt, error) {
	// TODO
	// 1. receive the request body image
	// 2. iteraction with gemini to extraction data
	// 3. get result the extraction
	// 4. save to db
	// 5. send to client

	r.Llm.Context = ctx // replace with request context
	receipt, err := r.Llm.ScanReceiptFromOCR(rawOcr)
	if err != nil {
		return nil, err
	}

	if !receipt.AIResponse.Success {
		return nil, responses.Err(http.StatusBadRequest, receipt.AIResponse.Message)
	}

	return nil, nil
}

func (r *ReceiptService) Scan(ctx context.Context, image []byte) (*models.Receipt, error) {
	// TODO
	// - receive the request body image
	// - iteraction with gemini to extraction data
	// - get result the extraction
	// - Check duplicate receipt
	// - save to db
	// - send to client

	r.Llm.Context = ctx // replace with request context
	receipt, err := r.Llm.ScanReceiptWithImage(image)
	if err != nil {
		return nil, err
	}

	if !receipt.AIResponse.Success {
		return nil, responses.BodyErr(receipt.AIResponse.Message)
	}

	// TODO check duplicate

	model := receipt.Model()
	res, err := r.ReceiptRepository.Save(ctx, model)
	if err != nil {
		logger.Log.Service(ctx).WithField("ai_response", receipt).WithField("error", err).Error("error when save receipt with success extraction data")
		return nil, err
	}

	return res, nil
}
