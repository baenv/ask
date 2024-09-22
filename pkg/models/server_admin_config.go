package models

import (
	"time"

	"gorm.io/gorm"
)

// ServerAdminConfig represents a server's admin configuration
type ServerAdminConfig struct {
	ID          int64     `json:"id" db:"id"`
	ServerID    int64     `json:"server_id" db:"server_id"`
	APIKey      string    `json:"api_key" db:"api_key"`
	EndpointURL string    `json:"endpoint_url" db:"endpoint_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	Command     string    `json:"command" db:"command"`
	Description string    `json:"description" db:"description"`

	Server *Server `json:"server" db:"-"`
}

// BeforeCreate is a GORM hook that generates a unique ID for the ServerAdminConfig
// and deactivates all other configs for the same server before activating this one
func (sac *ServerAdminConfig) BeforeCreate(tx *gorm.DB) error {
	if sac.ID == 0 {
		sac.ID = serverAdminConfigIDGenerator.Generate().Int64()
	}

	sac.IsActive = true

	return nil
}
