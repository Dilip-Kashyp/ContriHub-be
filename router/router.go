package router

import (
	"contrihub/constants"
	"contrihub/handlers"
	"contrihub/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	api := r.Group(constants.ApiV1)
	{
		// Auth routes
		api.GET(constants.AuthLogin, handlers.LoginHandler)

		// GitHub Proxy
		api.Any(constants.GithubProxy, handlers.ProxyHandler)

		// AI routes (rate-limited: 10 requests per 60 seconds per token)
		aiLimiter := middleware.NewRateLimiter(10, 60)
		aiGroup := api.Group("", aiLimiter.Middleware())
		{
			aiGroup.POST(constants.AIExplainRepo, handlers.ExplainRepoHandler)
			aiGroup.POST(constants.AIFindProjects, handlers.FindProjectsHandler)
			aiGroup.POST(constants.AIRoadmap, handlers.RoadmapHandler)
			aiGroup.POST(constants.AIStartGuide, handlers.StartGuideHandler)
			aiGroup.POST(constants.AIGenerateReadme, handlers.GenerateReadmeHandler)
			aiGroup.POST(constants.AIGenerateSummary, handlers.GenerateSummaryHandler)
			aiGroup.GET(constants.AIChatHistory, handlers.GetChatHistoryHandler)
			aiGroup.POST(constants.AIChatMessage, handlers.SubmitChatMessageHandler)
		}
	}

	// The callback doesn't have an api prefix based on the environment variable
	// CALLBACK_URL=http://localhost:8080/auth/github/callback
	// We'll set it at the root
	r.GET(constants.AuthCallback, handlers.CallbackHandler)

	return r
}
