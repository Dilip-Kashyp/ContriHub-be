package models

import "time"

// AICache stores LLM responses to avoid duplicate calls.
type AICache struct {
	ID           uint      `gorm:"primaryKey"`
	Key          string    `gorm:"uniqueIndex;size:64;not null"` // SHA-256 hash
	InputJSON    string    `gorm:"type:text;not null"`
	ResponseText string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (AICache) TableName() string {
	return "ai_cache"
}

// AIChatMessage stores chat history for a user (identified by token hash)
type AIChatMessage struct {
	ID        uint      `gorm:"primaryKey"`
	TokenHash string    `gorm:"index;size:64;not null"` // Identifier for the user
	Role      string    `gorm:"size:20;not null"`       // "user" or "ai"
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

func (AIChatMessage) TableName() string {
	return "ai_chat_messages"
}
