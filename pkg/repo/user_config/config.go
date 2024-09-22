package userconfig

import "gorm.io/gorm"

type userConfig struct {
	db *gorm.DB
}

func New(db *gorm.DB) IUserConfig {
	return &userConfig{db: db}
}
