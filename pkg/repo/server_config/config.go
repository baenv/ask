package serverconfig

import "gorm.io/gorm"

type serverConfig struct {
	db *gorm.DB
}

func New(db *gorm.DB) IServerConfig {
	return &serverConfig{db: db}
}
