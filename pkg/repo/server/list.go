package server

import "sum/pkg/models"

// ListByUserID
func (s server) ListByUserID(userID int64) ([]models.Server, error) {
	var servers []models.Server
	return servers, s.db.Where("owner_id = ?", userID).Find(&servers).Error
}
