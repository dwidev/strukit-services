package hash

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
	"strukit-services/internal/models"

	"github.com/google/uuid"
)

type ReceiptHashData struct {
	ProjectID       uuid.UUID
	MerchantName    *string
	ReceiptNumber   *string
	TotalAmount     float64
	TransactionDate string
	TransactionTime string
	Items           []*models.ReceiptItem
}

// GenerateContentHash creates a SHA-256 hash from receipt content
// This is used for exact duplicate detection
func GenerateContentHash(data ReceiptHashData) string {
	var parts []string

	if data.MerchantName != nil {
		parts = append(parts, strings.ToLower(strings.TrimSpace(*data.MerchantName)))
	} else {
		parts = append(parts, "")
	}

	if data.ReceiptNumber != nil {
		parts = append(parts, strings.ToLower(strings.TrimSpace(*data.ReceiptNumber)))
	} else {
		parts = append(parts, "")
	}

	// Add total amount (formatted to 2 decimal places)
	parts = append(parts, fmt.Sprintf("%.2f", data.TotalAmount))

	parts = append(parts, data.TransactionDate)
	parts = append(parts, data.TransactionTime)

	sortedItems := make([]*models.ReceiptItem, len(data.Items))
	copy(sortedItems, data.Items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].ItemName < sortedItems[j].ItemName
	})

	for _, item := range sortedItems {
		itemStr := fmt.Sprintf("%s:%d:%.2f",
			strings.ToLower(strings.TrimSpace(item.ItemName)),
			item.Quantity,
			item.TotalPrice)
		parts = append(parts, itemStr)
	}

	content := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}

func GenerateFingerprint(projectID uuid.UUID, contentHash string) string {
	content := fmt.Sprintf("%s:%s", projectID.String(), contentHash)
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}
