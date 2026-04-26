package constants

const (
	ErrMissingEnv   = "Missing required environment variables"
	ErrAuthFailed   = "Authentication failed"
	ErrProxyFailed  = "Failed to proxy request to GitHub"
	ErrNoCode       = "No authorization code provided"
	MsgLoginSuccess = "Login successful"

	// AI Standardized Codes
	ErrCodeRateLimit    = "rate_limit_exceeded"
	ErrCodeUnauthorized = "unauthorized"
	ErrCodeNotFound     = "not_found"
	ErrCodeServer       = "server_error"

	// AI User Messages
	ErrMsgRateLimit    = "Gibo is busy with too many requests. Please try again in a minute."
	ErrMsgUnauthorized = "Your session has expired. Please login again."
	ErrMsgNotFound     = "The requested repository or resource was not found."
	ErrMsgServer       = "Gibo is having some technical difficulties. Please try again later."
)
