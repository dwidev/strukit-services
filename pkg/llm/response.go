package llm

import (
	"strukit-services/internal/models"
	"strukit-services/pkg/helper"

	"github.com/google/uuid"
)

type ItemResponse struct {
	Name      *string `json:"name"`
	Quantity  *int    `json:"quantity"`
	UnitPrice *int    `json:"unitPrice"`
	Discount  *int    `json:"discount"`
	Total     *int    `json:"total"`
}

type PaymentSummaryResponse struct {
	PaymentMethod *string `json:"paymentMethod"` // TUNAI, CASH, QRIS, CREDIT_CARD, DEBIT, E_WALLET
	SubTotal      *int    `json:"subTotal"`
	Discount      *int    `json:"discount"`
	Tax           *int    `json:"tax"`
	AmountPaid    *int    `json:"amountPaid"`
	Paid          *int    `json:"paid"`
	Change        *int    `json:"change"`
}

type AIResponse struct {
	Success  bool   `json:"success"`
	Accuracy int    `json:"accuracy"`
	Message  string `json:"message"`
}

type ReceiptResponse struct {
	ReceiptNo      *string                 `json:"receiptNo"`
	ShopName       *string                 `json:"shopName"`
	Category       *string                 `json:"category"`
	AddressShop    *string                 `json:"addressShop"`
	ContactShop    *string                 `json:"contactShop"`
	Date           *string                 `json:"date"`
	Time           *string                 `json:"time"`
	CashierName    *string                 `json:"cashierName"`
	Items          *[]ItemResponse         `json:"items"`
	PaymentSummary *PaymentSummaryResponse `json:"paymentSummary"`
	AIResponse     AIResponse              `json:"aiResponse"`
}

func (rr *ReceiptResponse) Model() *models.Receipt {
	date := helper.ParseToDate(rr.Date)
	time := helper.ParseTimeOnly(*rr.Time)

	receiptID := uuid.New()
	accuracy := float64(rr.AIResponse.Accuracy) / 10
	subTotal := helper.IntPtrToFloat64(rr.PaymentSummary.SubTotal)
	Discount := helper.IntPtrToFloat64(rr.PaymentSummary.Discount)
	Tax := helper.IntPtrToFloat64(rr.PaymentSummary.Tax)
	AmountPaid := helper.IntPtrToFloat64(rr.PaymentSummary.AmountPaid)
	Paid := helper.IntPtrToFloat64(rr.PaymentSummary.Paid)
	Change := helper.IntPtrToFloat64(rr.PaymentSummary.Change)

	var items []*models.ReceiptItem

	if len(*rr.Items) > 0 {
		for _, v := range *rr.Items {
			up := helper.IntPtrToFloat64(v.UnitPrice)
			tp := helper.IntPtrToFloat64(v.Total)
			da := helper.IntPtrToFloat64(v.Discount)

			i := &models.ReceiptItem{
				ID:             uuid.New(),
				ReceiptID:      receiptID,
				ItemName:       *v.Name,
				Quantity:       *v.Quantity,
				UnitPrice:      &up,
				TotalPrice:     tp,
				ItemCode:       nil,
				Category:       nil,
				DiscountAmount: da,
				LineNumber:     nil,
			}

			items = append(items, i)
		}
	}

	return &models.Receipt{
		BaseModel: models.BaseModel{
			ID: receiptID,
		},
		ReceiptNumber:        rr.ReceiptNo,
		MerchantName:         rr.ShopName,
		SubTotal:             subTotal,
		Discount:             Discount,
		Tax:                  &Tax,
		TotalAmount:          AmountPaid,
		Paid:                 Paid,
		Change:               Change,
		TransactionDate:      *date,
		TransactionTime:      (*models.OnlyTime)(time),
		PaymentMethod:        rr.PaymentSummary.PaymentMethod,
		IsVerified:           true,
		VerificationNotes:    nil,
		ExtractionConfidence: &accuracy,
		AIModelUsed:          "gemini",
		Fingerprint:          "",
		ContentHash:          "",
		Items:                items,
		Category: &models.Category{
			Name: *rr.Category,
		},
	}
}
