package models

import (
	"time"

	"github.com/google/uuid"
)

type ReceiptItem struct {
	ID        uuid.UUID `json:"id"`
	ReceiptID uuid.UUID `json:"receipt_id"`

	ItemName   string   `json:"item_name"`
	Quantity   int      `json:"quantity"`
	UnitPrice  *float64 `json:"unit_price,omitempty"`
	TotalPrice float64  `json:"total_price"`

	ItemCode       *string `json:"item_code,omitempty"`
	Category       *string `json:"category,omitempty"`
	DiscountAmount float64 `json:"discount_amount"`

	LineNumber *int      `json:"line_number,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
	