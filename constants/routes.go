package constants

const (
	ApiV1        = "/api/v1"
	AuthLogin    = "/auth/login"
	AuthCallback = "/auth/github/callback"
	GithubProxy  = "/github/*path"

	// AI routes
	AIExplainRepo    = "/ai/explain-repo"
	AIFindProjects   = "/ai/find-projects"
	AIRoadmap        = "/ai/roadmap"
	AIStartGuide     = "/ai/start-guide"
	AIGenerateReadme = "/ai/generate-readme"
	AIGenerateSummary = "/ai/generate-summary"
	AIChatHistory     = "/ai/chat/history"
	AIChatMessage     = "/ai/chat/message"
)
