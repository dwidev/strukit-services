package models

import (
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	BaseModel
	UserID    uuid.UUID `json:"user_id"`
	ProjectID uuid.UUID `json:"project_id,"`

	// Extracted Data
	ReceiptNumber   *string    `json:"receiptNumber"`
	MerchantName    *string    `json:"merchantName"`
	SubTotal        float64    `json:"subTotal"`
	Discount        float64    `json:"discount"`
	Tax             *float64   `json:"tax"`
	TotalAmount     float64    `json:"totalAmount"`
	Paid            float64    `json:"paid"`
	Change          float64    `json:"change"`
	TransactionDate time.Time  `json:"transactionDate"`
	TransactionTime *time.Time `json:"transactionTime"`
	PaymentMethod   *string    `json:"paymentMethod"`

	// User Verification
	IsVerified        bool    `json:"isVerified"`
	VerificationNotes *string `json:"verificationNotes"`

	// AI Processing
	ExtractionConfidence *float64 `json:"extractionConfidence"`
	AIModelUsed          string   `json:"ai_modelUsed"`

	// Duplication info
	Fingerprint string `json:"fingerprint"`
	ContentHash string `json:"contentHash"`

	Items []*ReceiptItem `gorm:"foreignKey:ReceiptID;" json:"items"`
}
