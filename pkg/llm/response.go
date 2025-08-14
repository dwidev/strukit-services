package llm

type ItemResponse struct {
	Name      *string `json:"name"`
	Quantity  *int    `json:"quantity"`
	UnitPrice *int    `json:"unitPrice"`
	Discount  *int    `json:"discount"`
	Total     *int    `json:"total"`
}

type PaymentSummaryResponse struct {
	PaymentType *string `json:"paymentType"` // TUNAI, CASH, QRIS, CREDIT_CARD, DEBIT, E_WALLET
	SubTotal    *int    `json:"subTotal"`
	Tax         *int    `json:"tax"`
	AmountPaid  *int    `json:"amountPaid"`
	Paid        *int    `json:"paid"`
	Change      *int    `json:"change"`
}

type AIResponse struct {
	Success  bool    `json:"success"`
	Accuracy *int    `json:"accuracy"`
	Message  *string `json:"message"`
}

type ReceiptResponse struct {
	ReceiptNo      *string                 `json:"receiptNo"`
	ShopName       *string                 `json:"shopName"`
	Category       *string                 `json:"category"`
	AddressShop    *string                 `json:"addressShop"`
	ContactShop    *string                 `json:"contactShop"`
	Date           *string                 `json:"date"`
	CashierName    *string                 `json:"cashierName"`
	Items          *[]ItemResponse         `json:"items"`
	PaymentSummary *PaymentSummaryResponse `json:"paymentSummary"`
	AIResponse     *AIResponse             `json:"aiResponse"`
}
