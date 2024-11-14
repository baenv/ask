package server

import "sum/pkg/models"

func (s server) GetByServerID(serverID int64) (models.Server, error) {
	var server models.Server
	return server, s.db.Where("server_id = ?", serverID).First(&server).Error
}

func (s server) GetByPlatformID(userID, platform string) (models.Server, error) {
	var server models.Server
	return server, s.db.Where("server_id = ? AND platform = ?", userID, platform).First(&server).Error
}
