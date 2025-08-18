package appContext

type contextKey string

const (
	// client
	RequestIDKey contextKey = "request_id"
	IPAddressKey contextKey = "ip_address"
	PathKey      contextKey = "path"
	MethodKey    contextKey = "method"

	// token usage
	UserIDKey contextKey = "user_id"

	// client data
	ProjectID    contextKey = "project_id"
	ReceiptIDKey contextKey = "receipt_id"
)
