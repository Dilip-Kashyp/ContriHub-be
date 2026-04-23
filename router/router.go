package router

import (
	"contrihub/constants"
	"contrihub/handlers"

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
	}

	// The callback doesn't have an api prefix based on the environment variable
	// CALLBACK_URL=http://localhost:8080/auth/github/callback
	// We'll set it at the root
	r.GET(constants.AuthCallback, handlers.CallbackHandler)

	return r
}
