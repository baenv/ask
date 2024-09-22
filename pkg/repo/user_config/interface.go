package userconfig

import "ask/pkg/models"

type IUserConfig interface {
	ListPendingByUserID(userID string) ([]models.UserAgentConfig, error)
	GetByID(id string) (models.UserAgentConfig, error)
	SaveAPIKey(id string, apiKey string) error
	SaveEndpointURL(id string, endpointURL string) error
	SaveCommand(id string, command string) error
	RemoveByID(id string) error
	GetActiveByUserPlatformID(userID, platform string) (models.UserAgentConfig, error)
	SaveDescription(id string, description string) error
	ListByUserID(userID int64) ([]models.UserAgentConfig, error)
	GetByUserIDAndCommand(userID int64, command string) (models.UserAgentConfig, error)
}
