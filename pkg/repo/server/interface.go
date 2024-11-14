package server

import "sum/pkg/models"

type IServer interface {
	Create(server models.Server) (models.Server, error)
	GetByServerID(serverID int64) (models.Server, error)
	ListByUserID(userID int64) ([]models.Server, error)
	GetByPlatformID(userID, platform string) (models.Server, error)
}
