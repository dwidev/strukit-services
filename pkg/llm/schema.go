package llm

import "google.golang.org/genai"

func (m Manager) StructuredOutput() *genai.Schema {
	structuredOutput := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"aiResponse": {
				Type:        genai.TypeObject,
				Description: "AI summary result for the extraction receipt",
				Properties: map[string]*genai.Schema{
					"success": {
						Type:        genai.TypeBoolean,
						Description: "for the make sure is extraction is correct, set false if the wrong extraction",
					},
					"accuracy": {
						Type:        genai.TypeInteger,
						Description: "the value of confidently of the result extraction, the valus is between 1 - 10",
					},
					"message": {
						Type:        genai.TypeString,
						Description: "if the success result is false generate message error with message Bahasa ",
					},
				},
			},
			"receiptNo": {
				Type:        genai.TypeString,
				Description: "Unique receipt or transaction number printed on the receipt",
			},
			"shopName": {
				Type:        genai.TypeString,
				Description: "Name of the store or business establishment",
			},
			"category": {
				Type:        genai.TypeString,
				Description: "Category or type of receipt purpose (e.g., Belanja, Transport, Komsumsi)",
				Example:     "Belanja, Transport, Komsumsi",
			},
			"addressShop": {
				Type:        genai.TypeString,
				Description: "Physical address of the store or business location",
			},
			"contactShop": {
				Type:        genai.TypeString,
				Description: "Store contact information (phone number, email, or website)",
			},
			"date": {
				Type:        genai.TypeString,
				Description: "Transaction date and time in format DD/MM/YYYY HH:MM TZ",
				Example:     "15/02/2025 00:00 WIB",
			},
			"cashierName": {
				Type:        genai.TypeString,
				Description: "Name or ID of the cashier who processed the receipt transaction",
			},
			"items": {
				Type:        genai.TypeArray,
				Description: "List of purchased items with their details",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name": {
							Type:        genai.TypeString,
							Description: "Product or item name",
						},
						"quantity": {
							Type:        genai.TypeInteger,
							Description: "Quantity or amount of the item purchased",
						},
						"unitPrice": {
							Type:        genai.TypeInteger,
							Description: "Price per unit of the item",
						},
						"discount": {
							Type:        genai.TypeInteger,
							Description: "Discount amount applied to this item (if any)",
						},
						"total": {
							Type:        genai.TypeInteger,
							Description: "Total price for this item",
						},
					},
				},
			},
			"paymentSummary": {
				Type:        genai.TypeObject,
				Description: "Payment summary containing financial details of the receipt transaction",
				Properties: map[string]*genai.Schema{
					"paymentType": {
						Type:        genai.TypeString,
						Description: "Method of payment used (e.g., TUNTAI or CASH, QRIS, CREDIT CARD, DEBIT, E_WALLET)",
						Example:     "TUNAI or CASH, QRIS, CREDIT_CARD, DEBIT, E_WALLET",
					},
					"subTotal": {
						Type:        genai.TypeInteger,
						Description: "Subtotal amount before tax and other charges",
					},
					"tax": {
						Type:        genai.TypeInteger,
						Description: "Tax amount applied to the transaction",
					},
					"amountPaid": {
						Type:        genai.TypeInteger,
						Description: "Total amount that needs to be paid (subtotal + tax)",
					},
					"paid": {
						Type:        genai.TypeInteger,
						Description: "Actual amount paid by the customer",
					},
					"change": {
						Type:        genai.TypeInteger,
						Description: "Change amount returned to customer (paid - amountPaid)",
					},
				},
			},
		},
		PropertyOrdering: []string{
			"aiResponse",
			"receiptNo",
			"shopName",
			"category",
			"addressShop",
			"contactShop",
			"date",
			"cashierName",
			"items",
			"paymentSummary",
		},
		Required: []string{"aiResponse"},
	}
	return structuredOutput
}
