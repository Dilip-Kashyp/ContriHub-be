package handlers

import (
	"io"
	"net/http"
	"strings"

	"contrihub/constants"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	// Construct GitHub API URL
	githubApiUrl := "https://api.github.com" + path
	if len(c.Request.URL.RawQuery) > 0 {
		githubApiUrl += "?" + c.Request.URL.RawQuery
	}

	authHeader := c.GetHeader("Authorization")

	req, err := http.NewRequest(c.Request.Method, githubApiUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrProxyFailed})
		return
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	// Copy content type if body exists
	contentType := c.GetHeader("Content-Type")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": constants.ErrProxyFailed})
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for key, values := range resp.Header {
		// Avoid copying CORS headers from GitHub since our server handles CORS
		if strings.HasPrefix(strings.ToLower(key), "access-control-") {
			continue
		}
		for _, value := range values {
			c.Header(key, value)
		}
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
