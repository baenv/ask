package models

import (
	"time"

	"gorm.io/gorm"
)

// PlatformType represents the type of platform (Discord or Telegram)
type PlatformType string

const (
	PlatformDiscord  PlatformType = "discord"
	PlatformTelegram PlatformType = "telegram"
)

// User represents a user in the system
type User struct {
	ID         int64        `json:"id" db:"id"`
	UserID     string       `json:"user_id" db:"user_id"`
	Platform   PlatformType `json:"platform" db:"platform"`
	Username   string       `json:"username" db:"username"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	LastActive time.Time    `json:"last_active" db:"last_active"`

	UserAgentConfigs []UserAgentConfig `json:"user_agent_configs" gorm:"foreignKey:UserID"`
	Servers          []Server          `json:"servers" gorm:"foreignKey:OwnerID"`
}

// BeforeCreate(tx *gorm.DB) is a GORM hook that generates a unique ID for the User
// and handles conflicts based on user_id and platform
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == 0 {
		u.ID = userIDGenerator.Generate().Int64()
	}

	// Check for existing user with the same user_id and platform
	var existingUser User
	if tx.Where("user_id = ? AND platform = ?", u.UserID, u.Platform).First(&existingUser).Error == nil {
		u.ID = existingUser.ID
	}

	return nil
}
