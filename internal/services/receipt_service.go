package services

import (
	"context"
	"strukit-services/pkg/llm"
)

func NewReceipt(llm *llm.Manager) *ReceiptService {
	return &ReceiptService{
		Llm: llm,
	}
}

type ReceiptService struct {
	Llm *llm.Manager
}

func (r *ReceiptService) Scan(ctx context.Context, image []byte) (*string, error) {
	// func (r *ReceiptService) Scan(ctx context.Context) (*string, error) {
	// TODO
	// 1. receive the request body image
	// 2. iteraction with gemini to extraction data
	// 3. get result the extraction
	// 4. save to db
	// 5. send to client

	r.Llm.Context = ctx // replace with request context
	res, err := r.Llm.ScanReceiptWithImage(image)
	if err != nil {
		return nil, err
	}

	return res, nil
}
