package models

import (
	"time"

	"gorm.io/gorm"
)

// UserAgentConfig represents a user's agent configuration
type UserAgentConfig struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	APIKey      string    `json:"api_key" db:"api_key"`
	EndpointURL string    `json:"endpoint_url" db:"endpoint_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Command     string    `json:"command" db:"command"`
	Description string    `json:"description" db:"description"`
}

// BeforeCreate is a GORM hook that generates a unique ID for the UserAgentConfig
// and deactivates all other configs for the same user before activating this one
func (uac *UserAgentConfig) BeforeCreate(tx *gorm.DB) error {
	if uac.ID == 0 {
		uac.ID = userAgentConfigIDGenerator.Generate().Int64()
	}

	uac.IsActive = true

	return nil
}
