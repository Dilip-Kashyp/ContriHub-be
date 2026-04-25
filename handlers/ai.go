package handlers

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"contrihub/database"
	"contrihub/internal/ai"
	"contrihub/models"

	"github.com/gin-gonic/gin"
)

// ─── Request/Response Types ──────────────────────────────────────────────────

type ExplainRepoRequest struct {
	RepoName    string `json:"repo_name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Topics      string `json:"topics"`
	Readme      string `json:"readme"`
	Question    string `json:"question"`
}

type FindProjectsRequest struct {
	Query        string `json:"query"`
	RepoResults  string `json:"repo_results"` // Pre-formatted repo summaries from frontend
}

type RoadmapRequest struct {
	Interest   string `json:"interest"`
	SkillLevel string `json:"skill_level"`
	Repos      string `json:"repos"` // Pre-formatted repo summaries
}

type StartGuideRequest struct {
	RepoName      string `json:"repo_name"`
	Description   string `json:"description"`
	Language      string `json:"language"`
	Readme        string `json:"readme"`
	FileStructure string `json:"file_structure"`
}

type GenerateReadmeRequest struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	TopRepos  string `json:"top_repos"`
	Languages string `json:"languages"`
}

type GenerateSummaryRequest struct {
	Skills     string `json:"skills"`
	Projects   string `json:"projects"`
	Experience string `json:"experience"`
}

type AIResponse struct {
	Response string `json:"response"`
	Cached   bool   `json:"cached"`
}

type ChatMessageRequest struct {
	Message string `json:"message"`
}

type ChatMessageResponse struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ─── Handlers ────────────────────────────────────────────────────────────────

// ExplainRepoHandler handles POST /ai/explain-repo
func ExplainRepoHandler(c *gin.Context) {
	var req ExplainRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.RepoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo_name is required"})
		return
	}

	input := map[string]interface{}{
		"repo_name":   req.RepoName,
		"description": req.Description,
		"language":    req.Language,
		"question":    req.Question,
	}

	// Check cache
	cacheKey := ai.GenerateCacheKey("explain-repo", input)
	if cached, found := ai.GetCachedResponse(cacheKey); found {
		c.JSON(http.StatusOK, AIResponse{Response: cached, Cached: true})
		return
	}

	prompt := ai.BuildExplainRepoPrompt(req.RepoName, req.Description, req.Language, req.Topics, req.Readme, req.Question)
	response, err := ai.CallLLM(prompt, 1024)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	ai.SetCachedResponse(cacheKey, input, response)
	c.JSON(http.StatusOK, AIResponse{Response: response, Cached: false})
}

// FindProjectsHandler handles POST /ai/find-projects
func FindProjectsHandler(c *gin.Context) {
	var req FindProjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	input := map[string]interface{}{
		"query":       req.Query,
		"repo_results": req.RepoResults,
	}

	cacheKey := ai.GenerateCacheKey("find-projects", input)
	if cached, found := ai.GetCachedResponse(cacheKey); found {
		c.JSON(http.StatusOK, AIResponse{Response: cached, Cached: true})
		return
	}

	prompt := ai.BuildFindProjectsPrompt(req.Query, req.RepoResults)
	response, err := ai.CallLLM(prompt, 1024)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	ai.SetCachedResponse(cacheKey, input, response)
	c.JSON(http.StatusOK, AIResponse{Response: response, Cached: false})
}

// RoadmapHandler handles POST /ai/roadmap
func RoadmapHandler(c *gin.Context) {
	var req RoadmapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Interest == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interest is required"})
		return
	}

	input := map[string]interface{}{
		"interest":    req.Interest,
		"skill_level": req.SkillLevel,
	}

	cacheKey := ai.GenerateCacheKey("roadmap", input)
	if cached, found := ai.GetCachedResponse(cacheKey); found {
		c.JSON(http.StatusOK, AIResponse{Response: cached, Cached: true})
		return
	}

	prompt := ai.BuildRoadmapPrompt(req.Interest, req.SkillLevel, req.Repos)
	response, err := ai.CallLLM(prompt, 2048)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	ai.SetCachedResponse(cacheKey, input, response)
	c.JSON(http.StatusOK, AIResponse{Response: response, Cached: false})
}

// StartGuideHandler handles POST /ai/start-guide
func StartGuideHandler(c *gin.Context) {
	var req StartGuideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.RepoName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "repo_name is required"})
		return
	}

	input := map[string]interface{}{
		"repo_name": req.RepoName,
		"language":  req.Language,
	}

	cacheKey := ai.GenerateCacheKey("start-guide", input)
	if cached, found := ai.GetCachedResponse(cacheKey); found {
		c.JSON(http.StatusOK, AIResponse{Response: cached, Cached: true})
		return
	}

	prompt := ai.BuildStartGuidePrompt(req.RepoName, req.Description, req.Language, req.Readme, req.FileStructure)
	response, err := ai.CallLLM(prompt, 1536)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	ai.SetCachedResponse(cacheKey, input, response)
	c.JSON(http.StatusOK, AIResponse{Response: response, Cached: false})
}

// GenerateReadmeHandler handles POST /ai/generate-readme
func GenerateReadmeHandler(c *gin.Context) {
	var req GenerateReadmeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	input := map[string]interface{}{
		"username":  req.Username,
		"top_repos": req.TopRepos,
		"languages": req.Languages,
	}

	cacheKey := ai.GenerateCacheKey("generate-readme", input)
	if cached, found := ai.GetCachedResponse(cacheKey); found {
		c.JSON(http.StatusOK, AIResponse{Response: cached, Cached: true})
		return
	}

	prompt := ai.BuildReadmePrompt(req.Username, req.Name, req.Bio, req.TopRepos, req.Languages)
	response, err := ai.CallLLM(prompt, 2048)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	// Clean potential markdown code fences from response
	response = cleanMarkdownFences(response)

	ai.SetCachedResponse(cacheKey, input, response)
	c.JSON(http.StatusOK, AIResponse{Response: response, Cached: false})
}

// GenerateSummaryHandler handles POST /ai/generate-summary
func GenerateSummaryHandler(c *gin.Context) {
	var req GenerateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Skills == "" && req.Projects == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "skills or projects are required"})
		return
	}

	input := map[string]interface{}{
		"skills":   req.Skills,
		"projects": req.Projects,
	}

	cacheKey := ai.GenerateCacheKey("generate-summary", input)
	if cached, found := ai.GetCachedResponse(cacheKey); found {
		c.JSON(http.StatusOK, AIResponse{Response: cached, Cached: true})
		return
	}

	prompt := ai.BuildSummaryPrompt(req.Skills, req.Projects, req.Experience)
	response, err := ai.CallLLM(prompt, 512)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	ai.SetCachedResponse(cacheKey, input, response)
	c.JSON(http.StatusOK, AIResponse{Response: response, Cached: false})
}

// GetChatHistoryHandler retrieves the past 30 days of chat for the user
func GetChatHistoryHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	var messages []models.AIChatMessage
	if database.DB != nil {
		database.DB.Where("token_hash = ?", hash).Order("created_at asc").Find(&messages)
	}

	response := make([]ChatMessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = ChatMessageResponse{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	c.JSON(http.StatusOK, gin.H{"messages": response})
}

// SubmitChatMessageHandler processes a new user message and saves the interaction
func SubmitChatMessageHandler(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ChatMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))

	// Save user message
	if database.DB != nil {
		userMsg := models.AIChatMessage{
			TokenHash: hash,
			Role:      "user",
			Content:   req.Message,
		}
		database.DB.Create(&userMsg)
	}

	// Fetch recent history for context
	var history []models.AIChatMessage
	if database.DB != nil {
		database.DB.Where("token_hash = ?", hash).Order("created_at desc").Limit(10).Find(&history)
	}

	// Reverse to chronological
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	// Build context prompt
	var promptBuilder strings.Builder
	promptBuilder.WriteString("You are a helpful assistant for OpenSource developers on ContriHub.\n")
	promptBuilder.WriteString("You help developers find projects, understand codebases, and guide them.\n")
	promptBuilder.WriteString("Here is the recent conversation history:\n\n")

	for _, msg := range history {
		if msg.Role == "user" {
			promptBuilder.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		} else {
			promptBuilder.WriteString(fmt.Sprintf("AI: %s\n", msg.Content))
		}
	}
	promptBuilder.WriteString(fmt.Sprintf("User: %s\nAI: ", req.Message))

	// Call LLM
	response, err := ai.CallLLM(promptBuilder.String(), 1536)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("AI service error: %v", err)})
		return
	}

	response = cleanMarkdownFences(response)

	// Save AI message
	if database.DB != nil {
		aiMsg := models.AIChatMessage{
			TokenHash: hash,
			Role:      "ai",
			Content:   response,
		}
		database.DB.Create(&aiMsg)
	}

	c.JSON(http.StatusOK, ChatMessageResponse{
		Role:    "ai",
		Content: response,
	})
}

// cleanMarkdownFences removes wrapping ```markdown ... ``` fences from LLM output.
func cleanMarkdownFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```markdown") {
		s = strings.TrimPrefix(s, "```markdown")
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	} else if strings.HasPrefix(s, "```md") {
		s = strings.TrimPrefix(s, "```md")
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	} else if strings.HasPrefix(s, "```") {
		s = strings.TrimPrefix(s, "```")
		s = strings.TrimSuffix(s, "```")
		s = strings.TrimSpace(s)
	}
	return s
}
