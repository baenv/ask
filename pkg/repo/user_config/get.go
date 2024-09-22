package userconfig

import "ask/pkg/models"

func (c userConfig) GetByID(id string) (models.UserAgentConfig, error) {
	var config models.UserAgentConfig
	return config, c.db.Where("id = ?", id).First(&config).Error
}

func (c userConfig) GetActiveByUserPlatformID(userID, platform string) (models.UserAgentConfig, error) {
	var config models.UserAgentConfig
	return config, c.db.Table("user_agent_configs").
		Joins("JOIN users ON users.id = user_agent_configs.user_id").
		Where("users.user_id = ? AND users.platform = ? AND user_agent_configs.is_active = ?", userID, platform, true).
		First(&config).Error
}

func (c userConfig) GetByUserIDAndCommand(userID int64, command string) (models.UserAgentConfig, error) {
	var config models.UserAgentConfig
	return config, c.db.Table("user_agent_configs").
		Where("user_id = ? AND command = ?", userID, command).
		First(&config).Error
}
