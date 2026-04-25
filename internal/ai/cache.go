package ai

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"contrihub/database"
	"contrihub/models"
)

// GenerateCacheKey creates a deterministic SHA-256 hash from the endpoint and input data.
func GenerateCacheKey(endpoint string, input map[string]interface{}) string {
	// Sort keys for deterministic hashing
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	normalized := map[string]interface{}{"_endpoint": endpoint}
	for _, k := range keys {
		normalized[k] = input[k]
	}

	data, _ := json.Marshal(normalized)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// GetCachedResponse looks up a cached LLM response by key.
// Returns the response text and true if found, empty string and false otherwise.
func GetCachedResponse(key string) (string, bool) {
	if database.DB == nil {
		return "", false
	}

	var cache models.AICache
	result := database.DB.Where("key = ?", key).First(&cache)
	if result.Error != nil {
		return "", false
	}

	log.Printf("Cache HIT for key: %s", key[:16])
	return cache.ResponseText, true
}

// SetCachedResponse stores an LLM response in the cache.
func SetCachedResponse(key string, input map[string]interface{}, response string) {
	if database.DB == nil {
		return
	}

	inputJSON, _ := json.Marshal(input)
	cache := models.AICache{
		Key:          key,
		InputJSON:    string(inputJSON),
		ResponseText: response,
	}

	result := database.DB.Create(&cache)
	if result.Error != nil {
		log.Printf("Failed to cache response: %v", result.Error)
		return
	}

	log.Printf("Cache SET for key: %s", key[:16])
}
