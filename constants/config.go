package constants

import "os"

var (
	EnvClientId     = os.Getenv("CLIENT_ID")
	EnvClientSecret = os.Getenv("CLIENT_SECRET")
	EnvCallbackUrl  = os.Getenv("CALLBACK_URL")
	EnvFrontendUrl  = os.Getenv("FRONTEND_URL")
	EnvClientPat    = os.Getenv("CLIENT_PAT")
	EnvLLMApiKey    = os.Getenv("LLM_API_KEY")
	EnvLLMBaseUrl   = os.Getenv("LLM_BASE_URL")
)
