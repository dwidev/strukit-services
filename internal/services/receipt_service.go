package services

import (
	"context"
	"net/http"
	"strukit-services/pkg/llm"
	"strukit-services/pkg/responses"
)

func NewReceipt(llm *llm.Manager) *ReceiptService {
	return &ReceiptService{
		Llm: llm,
	}
}

type ReceiptService struct {
	Llm *llm.Manager
}

func (r *ReceiptService) ScanFromOCR(ctx context.Context, rawOcr string) (*llm.ReceiptResponse, error) {
	// func (r *ReceiptService) Scan(ctx context.Context) (*string, error) {
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

	return receipt, nil
}

func (r *ReceiptService) Scan(ctx context.Context, image []byte) (*llm.ReceiptResponse, error) {
	// func (r *ReceiptService) Scan(ctx context.Context) (*string, error) {
	// TODO
	// 1. receive the request body image
	// 2. iteraction with gemini to extraction data
	// 3. get result the extraction
	// 4. save to db
	// 5. send to client

	r.Llm.Context = ctx // replace with request context
	receipt, err := r.Llm.ScanReceiptWithImage(image)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
