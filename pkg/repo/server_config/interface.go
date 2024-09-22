package serverconfig

import "ask/pkg/models"

type IServerConfig interface {
	ListPendingByUserID(userID string) ([]models.ServerAdminConfig, error)
	GetByID(id string) (models.ServerAdminConfig, error)
	SaveAPIKey(id string, apiKey string) error
	SaveEndpointURL(id string, endpointURL string) error
	SaveCommand(id string, command string) error
	RemoveByID(id string) error
	GetActiveByServerPlatformID(serverID, platform string) (models.ServerAdminConfig, error)
	SaveDescription(id string, description string) error
	ListByServerID(serverID int64) ([]models.ServerAdminConfig, error)
	GetByServerIDAndCommand(serverID int64, command string) (models.ServerAdminConfig, error)
}
