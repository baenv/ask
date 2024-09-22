package userconfig

import "ask/pkg/models"

func (c userConfig) ListPendingByUserID(userID string) ([]models.UserAgentConfig, error) {
	var configs []models.UserAgentConfig
	return configs, c.db.Table("user_agent_configs").
		Select("user_agent_configs.*").
		Joins("JOIN users ON users.id = user_agent_configs.user_id").
		Where("users.user_id = ? AND (user_agent_configs.api_key = '' OR user_agent_configs.endpoint_url = '') OR user_agent_configs.is_active = false", userID).
		Find(&configs).Error
}

func (c userConfig) ListByUserID(userID int64) ([]models.UserAgentConfig, error) {
	var configs []models.UserAgentConfig
	return configs, c.db.Table("user_agent_configs").
		Where("user_id = ?", userID).
		Find(&configs).Error
}
