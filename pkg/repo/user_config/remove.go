package userconfig

import "sum/pkg/models"

func (c userConfig) RemoveByID(id string) error {
	return c.db.Delete(&models.UserAgentConfig{}, "id = ?", id).Error
}
