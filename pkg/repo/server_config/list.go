package serverconfig

import "ask/pkg/models"

func (c serverConfig) ListPendingByUserID(userID string) ([]models.ServerAdminConfig, error) {
	var configs []models.ServerAdminConfig
	return configs, c.db.Table("server_admin_configs").
		Preload("Server").
		Select("server_admin_configs.*").
		Joins("JOIN servers ON servers.id = server_admin_configs.server_id").
		Joins("JOIN users ON users.id = servers.owner_id").
		Where("users.user_id = ? AND (server_admin_configs.endpoint_url = '' OR server_admin_configs.api_key = '' OR server_admin_configs.command = '')", userID).
		Find(&configs).Error
}

func (c serverConfig) ListByServerID(serverID int64) ([]models.ServerAdminConfig, error) {
	var configs []models.ServerAdminConfig
	return configs, c.db.Table("server_admin_configs").
		Where("server_id = ?", serverID).
		Find(&configs).Error
}
