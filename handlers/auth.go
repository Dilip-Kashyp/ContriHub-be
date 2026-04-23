package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"contrihub/constants"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	clientId := os.Getenv(constants.EnvClientId)
	redirectUri := os.Getenv(constants.EnvCallbackUrl)

	if clientId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrMissingEnv})
		return
	}

	githubAuthUrl := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user,user:email",
		clientId,
		url.QueryEscape(redirectUri),
	)

	c.Redirect(http.StatusFound, githubAuthUrl)
}

func CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	frontendUrl := os.Getenv(constants.EnvFrontendUrl)
	if frontendUrl == "" {
		frontendUrl = "http://localhost:3000"
	}

	if code == "" {
		c.Redirect(http.StatusFound, frontendUrl+"/login?error=no_code")
		return
	}

	clientId := os.Getenv(constants.EnvClientId)
	clientSecret := os.Getenv(constants.EnvClientSecret)

	requestBody, _ := json.Marshal(map[string]string{
		"client_id":     clientId,
		"client_secret": clientSecret,
		"code":          code,
	})

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestBody))
	if err != nil {
		c.Redirect(http.StatusFound, frontendUrl+"/login?error=server_error")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Redirect(http.StatusFound, frontendUrl+"/login?error=server_error")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.Redirect(http.StatusFound, frontendUrl+"/login?error=server_error")
		return
	}

	if _, ok := result["error"]; ok {
		c.Redirect(http.StatusFound, frontendUrl+"/login?error=auth_failed")
		return
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		c.Redirect(http.StatusFound, frontendUrl+"/login?error=auth_failed")
		return
	}

	// Redirect back to frontend profile page with token
	redirectUrl := fmt.Sprintf("%s/profile#token=%s", frontendUrl, accessToken)
	c.Redirect(http.StatusFound, redirectUrl)
}
