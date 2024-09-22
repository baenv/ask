package serverconfig

import "ask/pkg/models"

func (c serverConfig) RemoveByID(id string) error {
	return c.db.Delete(&models.ServerAdminConfig{}, "id = ?", id).Error
}
