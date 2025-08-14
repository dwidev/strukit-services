package llm

import (
	"context"
	"strukit-services/pkg/logger"

	"google.golang.org/genai"
)

func Run(ctx context.Context) (*Manager, error) {
	manager := new(Manager)
	manager.Context = ctx

	client, err := genai.NewClient(manager.Context, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	manager.client = client
	return manager, nil
}

type Manager struct {
	Context context.Context
	client  *genai.Client
}

func (m *Manager) ScanReceiptWithImage(image []byte) (*string, error) {
	contents := []*genai.Content{
		genai.NewContentFromBytes(image, "image/jpeg", genai.RoleUser),
	}

	ccfg := m.contentCfg()
	res, err := m.client.Models.GenerateContent(m.Context, "gemini-2.5-flash", contents, ccfg)

	if err != nil {
		logger.Log.LLM(m.Context).Errorf("llm error generate content, error : %s", err)
		return nil, err
	}
	result := res.Text()
	return &result, nil
}

func (m *Manager) contentCfg() *genai.GenerateContentConfig {
	sys := SystemPrompt()
	temp := float32(0.7)

	structuredOutput := m.StructuredOutput()
	return &genai.GenerateContentConfig{
		ResponseMIMEType:  "application/json",
		SystemInstruction: genai.NewContentFromText(sys, genai.RoleModel),
		Temperature:       &temp,
		ResponseSchema:    structuredOutput,
	}
}

func (m Manager) StructuredOutput() *genai.Schema {
	structuredOutput := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
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
				Description: "Name or ID of the cashier who processed the transaction",
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
							Type:        genai.TypeString,
							Description: "Quantity or amount of the item purchased",
						},
						"unitPrice": {
							Type:        genai.TypeString,
							Description: "Price per unit of the item",
						},
						"discount": {
							Type:        genai.TypeString,
							Description: "Discount amount applied to this item (if any)",
						},
						"total": {
							Type:        genai.TypeString,
							Description: "Total price for this item (quantity × unitPrice - discount)",
						},
					},
				},
			},
			"paymentSummary": {
				Type:        genai.TypeArray,
				Description: "Payment summary containing financial details of the transaction",
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"paymentType": {
							Type:        genai.TypeString,
							Description: "Method of payment used (e.g., TUNTAI or CASH, QRIS, CREDIT CARD, DEBIT, E_WALLET)",
							Example:     "TUNAI or CASH, QRIS, CREDIT_CARD, DEBIT, E_WALLET",
						},
						"subTotal": {
							Type:        genai.TypeString,
							Description: "Subtotal amount before tax and other charges",
						},
						"tax": {
							Type:        genai.TypeString,
							Description: "Tax amount applied to the transaction",
						},
						"amountPaid": {
							Type:        genai.TypeString,
							Description: "Total amount that needs to be paid (subtotal + tax)",
						},
						"paid": {
							Type:        genai.TypeString,
							Description: "Actual amount paid by the customer",
						},
						"change": {
							Type:        genai.TypeString,
							Description: "Change amount returned to customer (paid - amountPaid)",
						},
					},
				},
			},
			"note": {
				Type:        genai.TypeString,
				Description: "Additional notes, terms, or information printed on the receipt",
			},
		},
		PropertyOrdering: []string{"shopName", "category", "addressShop"},
	}
	return structuredOutput
}
