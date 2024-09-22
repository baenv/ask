package models

import (
	"time"

	"gorm.io/gorm"
)

// Server represents a server (Discord) or group (Telegram) in the system
type Server struct {
	ID         int64        `json:"id" db:"id"`
	ServerID   string       `json:"server_id" db:"server_id"`
	OwnerID    string       `json:"owner_id" db:"owner_id"`
	Platform   PlatformType `json:"platform" db:"platform"`
	ServerName string       `json:"server_name" db:"server_name"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`

	ServerAdminConfig []ServerAdminConfig `json:"server_admin_config" gorm:"foreignKey:ServerID"`
}

// BeforeCreate(tx *gorm.DB) is a GORM hook that generates a unique ID for the Server
// and handles conflicts based on server_id and platform
func (s *Server) BeforeCreate(tx *gorm.DB) error {
	if s.ID == 0 {
		s.ID = serverIDGenerator.Generate().Int64()
	}

	// Check for existing server with the same server_id and platform
	var existingServer Server
	if tx.Where("server_id = ? AND platform = ?", s.ServerID, s.Platform).First(&existingServer).Error == nil {
		s.ID = existingServer.ID
	}

	return nil
}
