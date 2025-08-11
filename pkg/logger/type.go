package logger

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

const (
	CategoryRequest  = "REQUEST"
	CategoryDB       = "DATABASE"
	CategoryAI       = "AI_PROCESSING"
	CategoryBusiness = "BUSINESS_LOGIC"
	CategorySecurity = "SECURITY"
	CategoryPerf     = "PERFORMANCE"
	CategorySystem   = "SYSTEM"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	SessionIDKey contextKey = "session_id"
	IPAddressKey contextKey = "ip_address"
)
