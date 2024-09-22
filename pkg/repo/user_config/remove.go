package userconfig

import "ask/pkg/models"

func (c userConfig) RemoveByID(id string) error {
	return c.db.Delete(&models.UserAgentConfig{}, "id = ?", id).Error
}
