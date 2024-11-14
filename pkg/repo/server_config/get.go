package serverconfig

import "sum/pkg/models"

func (c serverConfig) GetByID(id string) (models.ServerAdminConfig, error) {
	var config models.ServerAdminConfig
	return config, c.db.Where("id = ?", id).Preload("Server").First(&config).Error
}

func (c serverConfig) GetActiveByServerPlatformID(serverID, platform string) (models.ServerAdminConfig, error) {
	var config models.ServerAdminConfig
	return config, c.db.Joins("JOIN servers ON servers.id = server_admin_configs.server_id").
		Where("servers.server_id = ? AND servers.platform = ? AND server_admin_configs.is_active = ?", serverID, platform, true).
		First(&config).Error
}

func (c serverConfig) GetByServerIDAndCommand(serverID int64, command string) (models.ServerAdminConfig, error) {
	var config models.ServerAdminConfig
	return config, c.db.Table("server_admin_configs").
		Where("server_id = ? AND command = ?", serverID, command).
		First(&config).Error
}
